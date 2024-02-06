package servicecheck

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Checker implements the kubenurse checker
type Checker struct {
	// Ingress and service config
	KubenurseIngressURL string
	KubenurseServiceURL string
	SkipCheckMeIngress  bool
	SkipCheckMeService  bool

	// shutdownDuration defines the time during which kubenurse will wait before stopping
	ShutdownDuration time.Duration

	// Kubernetes API
	KubernetesServiceHost    string
	KubernetesServicePort    string
	SkipCheckAPIServerDirect bool
	SkipCheckAPIServerDNS    bool

	// Neighbourhood
	KubenurseNamespace     string
	NeighbourFilter        string
	allowUnschedulable     bool
	SkipCheckNeighbourhood bool

	// TLS
	UseTLS bool

	// Controller runtime cached client
	client client.Client

	// metrics
	errorCounter      *prometheus.CounterVec
	durationHistogram *prometheus.HistogramVec

	// Http Client for https requests
	httpClient *http.Client

	// cachedResult represents a cached check result
	cachedResult *CachedResult

	// cacheTTL defines the TTL of how long a cached result is valid
	cacheTTL time.Duration

	// stop is used to cancel RunScheduled
	stop chan struct{}
}

// Result contains the result of a performed check run
type Result struct {
	APIServerDirect    string      `json:"api_server_direct"`
	APIServerDNS       string      `json:"api_server_dns"`
	MeIngress          string      `json:"me_ingress"`
	MeService          string      `json:"me_service"`
	NeighbourhoodState string      `json:"neighbourhood_state"`
	Neighbourhood      []Neighbour `json:"neighbourhood"`
}

// Check is the signature used by all checks that the checker can execute.
type Check func(ctx context.Context) (string, error)

// CachedResult represents a cached check result that is valid until the expiration.
type CachedResult struct {
	result     *Result
	expiration time.Time
}
