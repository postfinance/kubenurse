// Package servicecheck implements the checks the kubenurse performs.
package servicecheck

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	okStr            = "ok"
	errStr           = "error"
	skippedStr       = "skipped"
	MetricsNamespace = "kubenurse"
)

// New configures the checker with a httpClient and a cache timeout for check
// results. Other parameters of the Checker struct need to be configured separately.
func New(_ context.Context, cl client.Client, promRegistry *prometheus.Registry,
	allowUnschedulable bool, cacheTTL time.Duration, durationHistogramBuckets []float64) (*Checker, error) {
	// setup http transport
	tlsConfig, err := generateTLSConfig(os.Getenv("KUBENURSE_EXTRA_CA"))
	if err != nil {
		slog.Error("cannot generate tlsConfig with provided KUBENURSE_EXTRA_CA. Continuing with default tlsConfig",
			"KUBENURSE_EXTRA_CA", os.Getenv("KUBENURSE_EXTRA_CA"), "err", err)

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
		Transport: withHttptrace(promRegistry, transport, durationHistogramBuckets),
	}

	return &Checker{
		allowUnschedulable: allowUnschedulable,
		client:             cl,
		httpClient:         httpClient,
		cacheTTL:           cacheTTL,
		stop:               make(chan struct{}),
	}, nil
}

// Run runs all servicechecks and returns the result togeter with a boolean which indicates success. The cache
// is respected.
func (c *Checker) Run() {
	// Run Checks
	result := make(map[string]any)

	result[APIServerDirect] = c.measure(c.APIServerDirect, APIServerDirect)
	result[APIServerDNS] = c.measure(c.APIServerDNS, APIServerDNS)
	result[meIngress] = c.measure(c.MeIngress, meIngress)
	result[meService] = c.measure(c.MeService, meService)

	if c.SkipCheckNeighbourhood {
		result[NeighbourhoodState] = skippedStr
	} else {
		neighbours, err := c.GetNeighbours(context.Background(), c.KubenurseNamespace, c.NeighbourFilter)

		if err != nil {
			result[NeighbourhoodState] = err.Error()
		} else {
			result[NeighbourhoodState] = okStr
			result["neighbourhood"] = neighbours

			c.checkNeighbours(neighbours)
		}
	}

	// Cache result (used for /alive handler)
	c.LastCheckResult = result
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

	return c.doRequest(ctx, apiurl, false)
}

// APIServerDNS checks the /version endpoint of the Kubernetes API Server through the Cluster DNS URL
func (c *Checker) APIServerDNS(ctx context.Context) (string, error) {
	if c.SkipCheckAPIServerDNS {
		return skippedStr, nil
	}

	apiurl := fmt.Sprintf("https://kubernetes.default.svc.cluster.local:%s/version", c.KubernetesServicePort)

	return c.doRequest(ctx, apiurl, false)
}

// MeIngress checks if the kubenurse is reachable at the /alwayshappy endpoint behind the ingress
func (c *Checker) MeIngress(ctx context.Context) (string, error) {
	if c.SkipCheckMeIngress {
		return skippedStr, nil
	}

	return c.doRequest(ctx, c.KubenurseIngressURL+"/alwayshappy", false) //nolint:goconst // readability
}

// MeService checks if the kubenurse is reachable at the /alwayshappy endpoint through the kubernetes service
func (c *Checker) MeService(ctx context.Context) (string, error) {
	if c.SkipCheckMeService {
		return skippedStr, nil
	}

	return c.doRequest(ctx, c.KubenurseServiceURL+"/alwayshappy", false)
}

// measure implements metric collections for the check
func (c *Checker) measure(check Check, requestType string) string {
	// Add our label (check type) to the context so our http tracer can annotate
	// metrics and errors based with the label
	ctx := context.WithValue(context.Background(), kubenurseTypeKey{}, requestType)

	// Execute check
	res, _ := check(ctx) // this error is ignored as it is already logged in httptrace

	return res
}
