// Package servicecheck implements the checks the kubenurse performs.
package servicecheck

import (
	"context"
	"fmt"
	"log"
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
func New(ctx context.Context, discovery *kubediscovery.Client, promRegistry *prometheus.Registry,
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

	promRegistry.MustRegister(errorCounter, durationHistogram)

	// setup http transport
	transport, err := generateRoundTripper(os.Getenv("KUBENURSE_EXTRA_CA"), os.Getenv("KUBENURSE_INSECURE") == "true")
	if err != nil {
		log.Printf("using default transport: %s", err)

		transport = http.DefaultTransport
	}

	httpClient := &http.Client{
		Timeout:   5 * time.Second,
		Transport: withRequestTracing(promRegistry, transport),
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
func (c *Checker) APIServerDirect() (string, error) {
	if c.SkipCheckAPIServerDirect {
		return skippedStr, nil
	}

	apiurl := fmt.Sprintf("https://%s:%s/version", c.KubernetesServiceHost, c.KubernetesServicePort)

	return c.doRequest(apiurl)
}

// APIServerDNS checks the /version endpoint of the Kubernetes API Server through the Cluster DNS URL
func (c *Checker) APIServerDNS() (string, error) {
	if c.SkipCheckAPIServerDNS {
		return skippedStr, nil
	}

	apiurl := fmt.Sprintf("https://kubernetes.default.svc.cluster.local:%s/version", c.KubernetesServicePort)

	return c.doRequest(apiurl)
}

// MeIngress checks if the kubenurse is reachable at the /alwayshappy endpoint behind the ingress
func (c *Checker) MeIngress() (string, error) {
	if c.SkipCheckMeIngress {
		return skippedStr, nil
	}

	return c.doRequest(c.KubenurseIngressURL + "/alwayshappy") //nolint:goconst // readability
}

// MeService checks if the kubenurse is reachable at the /alwayshappy endpoint through the kubernetes service
func (c *Checker) MeService() (string, error) {
	if c.SkipCheckMeService {
		return skippedStr, nil
	}

	return c.doRequest(c.KubenurseServiceURL + "/alwayshappy")
}

// checkNeighbours checks the /alwayshappy endpoint from every discovered kubenurse neighbour. Neighbour pods on nodes
// which are not schedulable are excluded from this check to avoid possible false errors.
func (c *Checker) checkNeighbours(nh []kubediscovery.Neighbour) {
	for _, neighbour := range nh {
		neighbour := neighbour                 // pin
		if neighbour.Phase == v1.PodRunning && // only query running pods (excludes pending ones)
			!neighbour.Terminating && // exclude terminating pods
			(c.allowUnschedulable || neighbour.NodeSchedulable == kubediscovery.NodeSchedulable) {
			check := func() (string, error) {
				if c.UseTLS {
					return c.doRequest("https://" + neighbour.PodIP + ":8443/alwayshappy")
				}

				return c.doRequest("http://" + neighbour.PodIP + ":8080/alwayshappy")
			}

			_, _ = c.measure(check, "path_"+neighbour.NodeName)
		}
	}
}

// measure implements metric collections for the check
func (c *Checker) measure(check Check, label string) (string, error) {
	start := time.Now()

	// Execute check
	res, err := check()

	// Process metrics
	c.durationHistogram.WithLabelValues(label).Observe(time.Since(start).Seconds())

	if err != nil {
		log.Printf("failed request for %s with %v", label, err)
		c.errorCounter.WithLabelValues(label).Inc()
	}

	return res, err
}
