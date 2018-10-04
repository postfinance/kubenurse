package checker

import (
	"net/http"

	"github.com/postfinance/kubenurse/pkg/kubediscovery"
)

type Checker struct {
	// Ingress and service config
	KubenurseIngressUrl string
	KubenurseServiceUrl string

	// Kubernetes API
	KubernetesServiceHost string
	KubernetesServicePort string

	// Neighbourhood
	KubeNamespace   string
	NeighbourFilter string

	// Http Client for https requests
	HttpClient *http.Client
}

type Result struct {
	APIServerDirect    string                    `json:"api_server_direct"`
	APIServerDNS       string                    `json:"api_server_dns"`
	MeIngress          string                    `json:"me_ingress"`
	MeService          string                    `json:"me_service"`
	NeighbourhoodState string                    `json:"neighbourhood_state"`
	Neighbourhood      []kubediscovery.Neighbour `json:"neighbourhood"`
}

type Check func() (string, error)
