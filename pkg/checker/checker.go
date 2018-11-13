package checker

import (
	"fmt"
	"log"
	"time"

	"github.com/postfinance/kubenurse/pkg/kubediscovery"
	"github.com/postfinance/kubenurse/pkg/metrics"
)

// Run runs an check and returns the result
func (c *Checker) Run() (Result, bool) {
	var haserr bool
	var err error

	// Run Checks
	res := Result{}

	res.APIServerDirect, err = measure(c.ApiServerDirect, "api_server_direct")
	haserr = haserr || (err != nil)

	res.APIServerDNS, err = measure(c.ApiServerDNS, "api_server_dns")
	haserr = haserr || (err != nil)

	res.MeIngress, err = measure(c.MeIngress, "me_ingress")
	haserr = haserr || (err != nil)

	res.MeService, err = measure(c.MeService, "me_service")
	haserr = haserr || (err != nil)

	res.Neighbourhood, err = kubediscovery.GetNeighbourhood(c.KubenurseNamespace, c.NeighbourFilter)
	haserr = haserr || (err != nil)

	// Neighbourhood special error treating
	if err != nil {
		res.NeighbourhoodState = err.Error()
	} else {
		res.NeighbourhoodState = "ok"

		// Check all neighbours if the neighbourhood was discovered
		c.checkNeighbours(res.Neighbourhood)
	}

	return res, haserr
}

// RunScheduled runs the check run in the specified interval which can be used
// to keep the metrics up-to-date
func (c *Checker) RunScheduled(d time.Duration) {
	for range time.Tick(d) {
		c.Run()
	}
}

// ApiServerDirect checks the /version endpoint of the Kubernetes API Server through the direct link
func (c *Checker) ApiServerDirect() (string, error) {
	apiurl := fmt.Sprintf("https://%s:%s/version", c.KubernetesServiceHost, c.KubernetesServicePort)
	return c.doRequest(apiurl)
}

// ApiServerDNS checks the /version endpoint of the Kubernetes API Server through the Cluster DNS URL
func (c *Checker) ApiServerDNS() (string, error) {
	apiurl := fmt.Sprintf("https://kubernetes.default.svc:%s/version", c.KubernetesServicePort)
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

// checkNeighbours checks every provided neighbour at the /alwayshappy endpoint
func (c *Checker) checkNeighbours(nh []kubediscovery.Neighbour) {
	for _, neighbour := range nh {
		check := func() (string, error) {
			return c.doRequest("http://" + neighbour.PodIP + ":8080/alwayshappy")
		}

		measure(check, "path_"+neighbour.NodeName)
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
