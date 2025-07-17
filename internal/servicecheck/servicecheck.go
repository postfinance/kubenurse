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
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	okStr       = "ok"
	errStr      = "error"
	skippedStr  = "skipped"
	dialTimeout = 5 * time.Second
)

// New configures the checker with a httpClient and a cache timeout for check
// results. Other parameters of the Checker struct need to be configured separately.
func New(cl client.Client, allowUnschedulable bool, cacheTTL time.Duration, histogramGetter func(s string) Histogram) (*Checker, error) {
	// setup http transport
	tlsConfig, err := generateTLSConfig(os.Getenv("KUBENURSE_EXTRA_CA"))
	if err != nil {
		if !testing.Testing() {
			slog.Error("cannot generate tlsConfig with provided KUBENURSE_EXTRA_CA. Continuing with default tlsConfig",
				"KUBENURSE_EXTRA_CA", os.Getenv("KUBENURSE_EXTRA_CA"), "err", err)
		}

		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	tlsConfig.InsecureSkipVerify = os.Getenv("KUBENURSE_INSECURE") == "true"
	dialer := &net.Dialer{
		Timeout:   dialTimeout,
		KeepAlive: dialTimeout,
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
		Timeout:   dialTimeout + time.Second,
		Transport: withHttptrace(transport, histogramGetter),
	}

	return &Checker{
		allowUnschedulable: allowUnschedulable,
		client:             cl,
		httpClient:         httpClient,
		cacheTTL:           cacheTTL,
		ExtraChecks:        make(map[string]string),
	}, nil
}

// Run runs all servicechecks and returns the result togeter with a boolean which indicates success. The cache
// is respected.
func (c *Checker) Run(ctx context.Context) {
	// Run Checks
	result := sync.Map{}

	wg := sync.WaitGroup{}

	// Cache result (used for /alive handler)
	defer func() {
		res := make(map[string]any)

		result.Range(func(key, value any) bool {
			k, _ := key.(string)
			res[k] = value

			return true
		})

		c.LastCheckResult = res
	}()

	wg.Add(4)

	go c.measure(ctx, &wg, &result, c.APIServerDirect, APIServerDirect)
	go c.measure(ctx, &wg, &result, c.APIServerDNS, APIServerDNS)
	go c.measure(ctx, &wg, &result, c.MeIngress, meIngress)
	go c.measure(ctx, &wg, &result, c.MeService, meService)

	wg.Add(len(c.ExtraChecks))

	for metricName, url := range c.ExtraChecks {
		go c.measure(ctx, &wg, &result,
			func(ctx context.Context) string { return c.doRequest(ctx, url, false) },
			metricName)
	}

	if c.SkipCheckNeighbourhood {
		result.Store(NeighbourhoodState, skippedStr)
		return
	}

	neighbours, err := c.getNeighbours(ctx, c.KubenurseNamespace, c.NeighbourFilter)
	if err != nil {
		result.Store(NeighbourhoodState, err.Error())
		return
	}

	result.Store(NeighbourhoodState, okStr)
	result.Store(Neighbourhood, neighbours)

	if c.NeighbourLimit > 0 && len(neighbours) > c.NeighbourLimit {
		neighbours = c.filterNeighbours(neighbours)
	}

	wg.Add((len(neighbours)))

	for _, neighbour := range neighbours {
		check := func(ctx context.Context) string {
			return c.doRequest(ctx, podIPtoURL(neighbour.PodIP, c.UseTLS), true)
		}

		go c.measure(ctx, &wg, &result, check, "path_"+neighbour.NodeName)
	}

	wg.Wait()
}

// APIServerDirect checks the /version endpoint of the Kubernetes API Server through the direct link
func (c *Checker) APIServerDirect(ctx context.Context) string {
	if c.SkipCheckAPIServerDirect {
		return skippedStr
	}

	apiurl := fmt.Sprintf("https://%s/version", net.JoinHostPort(c.KubernetesServiceHost, c.KubernetesServicePort))

	return c.doRequest(ctx, apiurl, false)
}

// APIServerDNS checks the /version endpoint of the Kubernetes API Server through the Cluster DNS URL
func (c *Checker) APIServerDNS(ctx context.Context) string {
	if c.SkipCheckAPIServerDNS {
		return skippedStr
	}

	apiurl := fmt.Sprintf("https://%s/version", net.JoinHostPort(c.KubernetesServiceDNS, c.KubernetesServicePort))

	return c.doRequest(ctx, apiurl, false)
}

// MeIngress checks if the kubenurse is reachable at the /alwayshappy endpoint behind the ingress
func (c *Checker) MeIngress(ctx context.Context) string {
	if c.SkipCheckMeIngress {
		return skippedStr
	}

	return c.doRequest(ctx, c.KubenurseIngressURL+"/alwayshappy", false) //nolint:goconst // readability
}

// MeService checks if the kubenurse is reachable at the /alwayshappy endpoint through the kubernetes service
func (c *Checker) MeService(ctx context.Context) string {
	if c.SkipCheckMeService {
		return skippedStr
	}

	return c.doRequest(ctx, c.KubenurseServiceURL+"/alwayshappy", false)
}

// measure implements metric collections for the check
func (c *Checker) measure(ctx context.Context, wg *sync.WaitGroup, res *sync.Map, check Check, requestType string) {
	// Add our label (check type) to the context so our http tracer can annotate
	// metrics and errors based with the label
	defer wg.Done()

	ctx = context.WithValue(ctx, kubenurseTypeKey{}, requestType)
	ctx = context.WithValue(ctx, kubenurseErrorAccountedKey{}, &atomic.Bool{})
	res.Store(requestType, check(ctx))
}

func podIPtoURL(podIP string, useTLS bool) string {
	if useTLS {
		return "https://" + net.JoinHostPort(podIP, "8443") + "/alwayshappy"
	}

	return "http://" + net.JoinHostPort(podIP, "8080") + "/alwayshappy"
}
