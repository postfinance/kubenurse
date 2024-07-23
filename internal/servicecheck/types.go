package servicecheck

import (
	"context"
	"net/http"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	NeighbourhoodState = "neighbourhood_state"
	Neighbourhood      = "neighbourhood"
	meService          = "me_service"
	meIngress          = "me_ingress"
	APIServerDirect    = "api_server_direct"
	APIServerDNS       = "api_server_dns"
)

// Checker implements the kubenurse checker
type Checker struct {
	// Ingress and service config
	KubenurseIngressURL string
	KubenurseServiceURL string
	SkipCheckMeIngress  bool
	SkipCheckMeService  bool

	// shutdownDuration defines the time during which kubenurse will accept https requests during shutdown
	ShutdownDuration time.Duration

	// Kubernetes API
	KubernetesServiceHost    string
	KubernetesServicePort    string
	SkipCheckAPIServerDirect bool
	SkipCheckAPIServerDNS    bool

	// Neighbourhood
	KubenurseNamespace     string
	NeighbourFilter        string
	NeighbourLimit         int
	allowUnschedulable     bool
	SkipCheckNeighbourhood bool

	// Additional endpoints
	ExtraChecks map[string]string

	// TLS
	UseTLS bool

	// Controller runtime cached client
	client client.Client

	// Http Client for https requests
	httpClient *http.Client

	// LastCheckResult represents a cached check result
	LastCheckResult map[string]any

	// cacheTTL defines the TTL of how long a cached result is valid
	cacheTTL time.Duration
}

// Check is the signature used by all checks that the checker can execute.
type Check func(ctx context.Context) string
