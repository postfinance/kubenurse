package checker

import (
	"fmt"
	"log"
	"time"

	"github.com/postfinance/kubenurse/pkg/kubediscovery"
	"github.com/postfinance/kubenurse/pkg/metrics"
)

func (c *Checker) Run() (Result, bool) {
	var haserr bool
	var err error

	// Run Checks
	res := Result{}

	res.APIServerDirect, err = meassure(c.ApiServerDirect, "api_server_direct")
	haserr = haserr || (err != nil)

	res.APIServerDNS, err = meassure(c.ApiServerDNS, "api_server_dns")
	haserr = haserr || (err != nil)

	res.MeIngress, err = meassure(c.MeIngress, "me_ingress")
	haserr = haserr || (err != nil)

	res.MeService, err = meassure(c.MeService, "me_service")
	haserr = haserr || (err != nil)

	res.Neighbourhood, err = kubediscovery.GetNeighbourhood(c.KubeNamespace, c.NeighbourFilter)
	haserr = haserr || (err != nil)

	// Neighbourhood special error treating
	if err != nil {
		res.NeighbourhoodState = err.Error()
	} else {
		res.NeighbourhoodState = "ok"

		// Check all neighbours
		c.checkNeighbours(res.Neighbourhood)
	}

	return res, haserr
}

func (c *Checker) RunScheduled(d time.Duration) {
	for range time.Tick(d) {
		c.Run()
	}
}

func (c *Checker) ApiServerDirect() (string, error) {
	apiurl := fmt.Sprintf("https://%s:%s/version", c.KubernetesServiceHost, c.KubernetesServicePort)
	return c.doRequest(apiurl)
}

func (c *Checker) ApiServerDNS() (string, error) {
	apiurl := fmt.Sprintf("https://kubernetes.default.svc:%s/version", c.KubernetesServicePort)
	return c.doRequest(apiurl)
}

func (c *Checker) MeIngress() (string, error) {
	return c.doRequest(c.KubenurseIngressUrl + "/alwayshappy")
}

func (c *Checker) MeService() (string, error) {
	return c.doRequest(c.KubenurseServiceUrl + "/alwayshappy")
}

func (c *Checker) checkNeighbours(nh []kubediscovery.Neighbour) {
	for _, neighbour := range nh {
		check := func() (string, error) {
			return c.doRequest("http://" + neighbour.PodIP + ":8080/alwayshappy")
		}

		meassure(check, "path_"+neighbour.NodeName)
	}
}

func meassure(check Check, label string) (string, error) {
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
