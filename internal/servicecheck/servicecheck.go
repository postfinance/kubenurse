// Package servicecheck implements the checks the kubenurse performs.
package servicecheck

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/postfinance/kubenurse/internal/kubediscovery"
	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
)

const (
	okStr            = "ok"
	errStr           = "error"
	skippedStr       = "skipped"
	metricsNamespace = "kubenurse"
)

// New configures the checker with a httpClient and a cache timeout for check
// results. Other parameters of the Checker struct need to be configured separately.
func New(_ context.Context, discovery *kubediscovery.Client, promRegistry *prometheus.Registry,
	allowUnschedulable bool, cacheTTL time.Duration, durationHistogramBuckets []float64) (*Checker, error) {
	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "errors_total",
			Help:      "Kubenurse error counter partitioned by error type",
		},
		[]string{"type"},
	)

	durationHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "request_duration",
			Help:      "Kubenurse request duration partitioned by target path",
			Buckets:   durationHistogramBuckets,
		},
		[]string{"type"},
	)

	// TODO: Add label for which request it was as this is not helpful in this current state
	// TODO: Do we want to have it also as summary?
	latencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "httpclient_trace_request_duration_seconds",
			Help:      "Latency histogram for requests from the kubenurse http client. Time in seconds since the start of the http request.",
			Buckets:   []float64{.0005, .005, .01, .025, .05, .1, .25, .5, 1}, // TODO: Which buckets are really needed?
		},
		[]string{"event", "type"},
	)

	promRegistry.MustRegister(errorCounter, durationHistogram, latencyVec)

	// setup http transport
	tlsConfig, err := generateTLSConfig(os.Getenv("KUBENURSE_EXTRA_CA"))
	if err != nil {
		log.Printf("cannot generate tlsConfig with KUBENURSE_EXTRA_CA: %s", err)

		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	tlsConfig.InsecureSkipVerify = os.Getenv("KUBENURSE_INSECURE") == "true"
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		TLSClientConfig:       tlsConfig,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		DisableKeepAlives:     os.Getenv("KUBENURSE_REUSE_CONNECTIONS") != "true",
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient := &http.Client{
		Timeout:   5 * time.Second,
		Transport: withHttptrace(promRegistry, transport, latencyVec),
	}

	return &Checker{
		allowUnschedulable: allowUnschedulable,
		discovery:          discovery,
		httpClient:         httpClient,
		cacheTTL:           cacheTTL,
		errorCounter:       errorCounter,
		durationHistogram:  durationHistogram,
		stop:               make(chan struct{}),
	}, nil
}

// Run runs all servicechecks and returns the result togeter with a boolean which indicates success. The cache
// is respected.
func (c *Checker) Run() (Result, bool) {
	var (
		haserr bool
		err    error
	)

	// Check if a result is cached and return it
	cacheRes := c.retrieveResultFromCache()
	if cacheRes != nil {
		return *cacheRes, false
	}

	// Run Checks
	res := Result{}

	res.APIServerDirect, err = c.measure(c.APIServerDirect, "api_server_direct")
	haserr = haserr || (err != nil)

	res.APIServerDNS, err = c.measure(c.APIServerDNS, "api_server_dns")
	haserr = haserr || (err != nil)

	res.MeIngress, err = c.measure(c.MeIngress, "me_ingress")
	haserr = haserr || (err != nil)

	res.MeService, err = c.measure(c.MeService, "me_service")
	haserr = haserr || (err != nil)

	if c.SkipCheckNeighbourhood {
		res.NeighbourhoodState = skippedStr
	} else {
		res.Neighbourhood, err = c.discovery.GetNeighbours(context.TODO(), c.KubenurseNamespace, c.NeighbourFilter)
		haserr = haserr || (err != nil)

		// Neighbourhood special error treating
		if err != nil {
			res.NeighbourhoodState = err.Error()
		} else {
			res.NeighbourhoodState = okStr

			// Check all neighbours if the neighbourhood was discovered
			c.checkNeighbours(res.Neighbourhood)
		}
	}

	// Cache result
	c.cacheResult(&res)

	return res, haserr
}

// RunScheduled runs the checks in the specified interval which can be used to keep the metrics up-to-date. This
// function does not return until StopScheduled is called.
func (c *Checker) RunScheduled(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Run()
		case <-c.stop:
			return
		}
	}
}

// StopScheduled is used to stop the scheduled run of checks.
func (c *Checker) StopScheduled() {
	close(c.stop)
}

// APIServerDirect checks the /version endpoint of the Kubernetes API Server through the direct link
func (c *Checker) APIServerDirect(ctx context.Context) (string, error) {
	if c.SkipCheckAPIServerDirect {
		return skippedStr, nil
	}

	apiurl := fmt.Sprintf("https://%s:%s/version", c.KubernetesServiceHost, c.KubernetesServicePort)

	return c.doRequest(ctx, apiurl)
}

// APIServerDNS checks the /version endpoint of the Kubernetes API Server through the Cluster DNS URL
func (c *Checker) APIServerDNS(ctx context.Context) (string, error) {
	if c.SkipCheckAPIServerDNS {
		return skippedStr, nil
	}

	apiurl := fmt.Sprintf("https://kubernetes.default.svc.cluster.local:%s/version", c.KubernetesServicePort)

	return c.doRequest(ctx, apiurl)
}

// MeIngress checks if the kubenurse is reachable at the /alwayshappy endpoint behind the ingress
func (c *Checker) MeIngress(ctx context.Context) (string, error) {
	if c.SkipCheckMeIngress {
		return skippedStr, nil
	}

	return c.doRequest(ctx, c.KubenurseIngressURL+"/alwayshappy") //nolint:goconst // readability
}

// MeService checks if the kubenurse is reachable at the /alwayshappy endpoint through the kubernetes service
func (c *Checker) MeService(ctx context.Context) (string, error) {
	if c.SkipCheckMeService {
		return skippedStr, nil
	}

	return c.doRequest(ctx, c.KubenurseServiceURL+"/alwayshappy")
}

// checkNeighbours checks the /alwayshappy endpoint from every discovered kubenurse neighbour. Neighbour pods on nodes
// which are not schedulable are excluded from this check to avoid possible false errors.
func (c *Checker) checkNeighbours(nh []kubediscovery.Neighbour) {
	for _, neighbour := range nh {
		neighbour := neighbour                 // pin
		if neighbour.Phase == v1.PodRunning && // only query running pods (excludes pending ones)
			!neighbour.Terminating && // exclude terminating pods
			(c.allowUnschedulable || neighbour.NodeSchedulable == kubediscovery.NodeSchedulable) {
			check := func(ctx context.Context) (string, error) {
				if c.UseTLS {
					return c.doRequest(ctx, "https://"+neighbour.PodIP+":8443/alwayshappy")
				}

				return c.doRequest(ctx, "http://"+neighbour.PodIP+":8080/alwayshappy")
			}

			_, _ = c.measure(check, "path_"+neighbour.NodeName)
		}
	}
}

// measure implements metric collections for the check
func (c *Checker) measure(check Check, label string) (string, error) {
	start := time.Now()

	// Add our label (check type) to the context so our http tracer can annotate
	// metrics and errors based with the label
	ctx := context.WithValue(context.Background(), kubenurseContextKey{}, label)

	// Execute check
	res, err := check(ctx)

	// Process metrics
	c.durationHistogram.WithLabelValues(label).Observe(time.Since(start).Seconds())

	if err != nil {
		log.Printf("failed request for %s with %v", label, err)
		c.errorCounter.WithLabelValues(label).Inc()
	}

	return res, err
}
