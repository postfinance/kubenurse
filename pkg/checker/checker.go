// Package checker implements the checks the kubenurse performs.
package checker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/postfinance/kubenurse/pkg/kubediscovery"
	"github.com/postfinance/kubenurse/pkg/metrics"
)

// New configures the checker with a httpClient and a cache timeout for check
// results. Other parameters of the Checker struct need to be configured separately.
func New(ctx context.Context, httpClient *http.Client, cacheTTL time.Duration, allowUnschedulable bool) (*Checker, error) {
	discovery, err := kubediscovery.New(ctx, allowUnschedulable)
	if err != nil {
		return nil, fmt.Errorf("create k8s discovery client: %w", err)
	}

	return &Checker{
		allowUnschedulable: allowUnschedulable,
		discovery:          discovery,
		httpClient:         httpClient,
		cacheTTL:           cacheTTL,
	}, nil
}

// Run runs an check and returns the result togeter with a boolean, if it wasn't
// successful. It respects the cache.
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

	res.APIServerDirect, err = measure(c.APIServerDirect, "api_server_direct")
	haserr = haserr || (err != nil)

	res.APIServerDNS, err = measure(c.APIServerDNS, "api_server_dns")
	haserr = haserr || (err != nil)

	res.MeIngress, err = measure(c.MeIngress, "me_ingress")
	haserr = haserr || (err != nil)

	res.MeService, err = measure(c.MeService, "me_service")
	haserr = haserr || (err != nil)

	res.Neighbourhood, err = c.discovery.GetNeighbours(context.TODO(), c.KubenurseNamespace, c.NeighbourFilter)
	haserr = haserr || (err != nil)

	// Neighbourhood special error treating
	if err != nil {
		res.NeighbourhoodState = err.Error()
	} else {
		res.NeighbourhoodState = "ok"

		// Check all neighbours if the neighbourhood was discovered
		c.checkNeighbours(res.Neighbourhood)
	}

	// Cache result
	c.cacheResult(&res)

	return res, haserr
}

// RunScheduled runs the check run in the specified interval which can be used
// to keep the metrics up-to-date.
func (c *Checker) RunScheduled(d time.Duration) {
	for range time.Tick(d) {
		c.Run()
	}
}

// APIServerDirect checks the /version endpoint of the Kubernetes API Server through the direct link
func (c *Checker) APIServerDirect() (string, error) {
	apiurl := fmt.Sprintf("https://%s:%s/version", c.KubernetesServiceHost, c.KubernetesServicePort)
	return c.doRequest(apiurl)
}

// APIServerDNS checks the /version endpoint of the Kubernetes API Server through the Cluster DNS URL
func (c *Checker) APIServerDNS() (string, error) {
	apiurl := fmt.Sprintf("https://kubernetes.default.svc.cluster.local:%s/version", c.KubernetesServicePort)
	return c.doRequest(apiurl)
}

// MeIngress checks if the kubenurse is reachable at the /alwayshappy endpoint behind the ingress
func (c *Checker) MeIngress() (string, error) {
	return c.doRequest(c.KubenurseIngressURL + "/alwayshappy")
}

// MeService checks if the kubenurse is reachable at the /alwayshappy endpoint through the kubernetes service
func (c *Checker) MeService() (string, error) {
	return c.doRequest(c.KubenurseServiceURL + "/alwayshappy")
}

// checkNeighbours checks the /alwayshappy endpoint from every discovered kubenurse neighbour. Neighbour pods on nodes
// which are not schedulable are excluded from this check to avoid possible false errors.
func (c *Checker) checkNeighbours(nh []kubediscovery.Neighbour) {
	for _, neighbour := range nh {
		neighbour := neighbour // pin
		if c.allowUnschedulable || neighbour.NodeSchedulable == kubediscovery.NodeSchedulable {
			check := func() (string, error) {
				if c.UseTLS {
					return c.doRequest("https://" + neighbour.PodIP + ":8443/alwayshappy")
				}

				return c.doRequest("http://" + neighbour.PodIP + ":8080/alwayshappy")
			}

			_, _ = measure(check, "path_"+neighbour.NodeName)
		}
	}
}

// measure implements metric collections for the check
func measure(check Check, label string) (string, error) {
	start := time.Now()

	// Execute check
	res, err := check()

	// Process metrics
	metrics.DurationSummary.WithLabelValues(label).Observe(time.Since(start).Seconds())

	if err != nil {
		log.Printf("failed request for %s with %v", label, err)
		metrics.ErrorCounter.WithLabelValues(label).Inc()
	}

	return res, err
}
