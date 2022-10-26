// Package kubenurse contains the server code for the kubenurse service.
package kubenurse

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/postfinance/kubenurse/internal/kubediscovery"
	"github.com/postfinance/kubenurse/internal/servicecheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
)

// Server is used to build the kubenurse http/https server(s).
type Server struct {
	http  http.Server
	https http.Server

	checker *servicecheck.Checker

	// Configuration options
	useTLS bool
	// If we want to consider kubenurses on unschedulable nodes
	allowUnschedulable bool

	// Mutex to protect ready flag
	mu    *sync.Mutex
	ready bool
}

// New creates a new kubenurse server. The server can be configured with the following environment variables:
// * KUBENURSE_USE_TLS
// * KUBENURSE_ALLOW_UNSCHEDULABLE
// * KUBENURSE_INGRESS_URL
// * KUBENURSE_SERVICE_URL
// * KUBERNETES_SERVICE_HOST
// * KUBERNETES_SERVICE_PORT
// * KUBENURSE_NAMESPACE
// * KUBENURSE_NEIGHBOUR_FILTER
// * KUBENURSE_CHECK_API_SERVER_DIRECT
// * KUBENURSE_CHECK_API_SERVER_DNS
// * KUBENURSE_CHECK_ME_INGRESS
// * KUBENURSE_CHECK_ME_SERVICE
// * KUBENURSE_CHECK_NEIGHBOURHOOD
func New(ctx context.Context, k8s kubernetes.Interface) (*Server, error) {
	mux := http.NewServeMux()

	server := &Server{
		http: http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
		https: http.Server{
			Addr:              ":8443",
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       120 * time.Second,
		},

		//nolint:goconst // No need to make "true" a constant in my opinion, readability is better like this.
		useTLS:             os.Getenv("KUBENURSE_USE_TLS") == "true",
		allowUnschedulable: os.Getenv("KUBENURSE_ALLOW_UNSCHEDULABLE") == "true",

		mu:    new(sync.Mutex),
		ready: true,
	}

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	discovery, err := kubediscovery.New(ctx, k8s, server.allowUnschedulable)
	if err != nil {
		return nil, fmt.Errorf("create k8s discovery client: %w", err)
	}

	// setup checker
	chk, err := servicecheck.New(ctx, discovery, promRegistry, server.allowUnschedulable, 3*time.Second)
	if err != nil {
		return nil, err
	}

	chk.KubenurseIngressURL = os.Getenv("KUBENURSE_INGRESS_URL")
	chk.KubenurseServiceURL = os.Getenv("KUBENURSE_SERVICE_URL")
	chk.KubernetesServiceHost = os.Getenv("KUBERNETES_SERVICE_HOST")
	chk.KubernetesServicePort = os.Getenv("KUBERNETES_SERVICE_PORT")
	chk.KubenurseNamespace = os.Getenv("KUBENURSE_NAMESPACE")
	chk.NeighbourFilter = os.Getenv("KUBENURSE_NEIGHBOUR_FILTER")

	//nolint:goconst // No need to make "false" a constant in my opinion, readability is better like this.
	chk.SkipCheckAPIServerDirect = os.Getenv("KUBENURSE_CHECK_API_SERVER_DIRECT") == "false"
	chk.SkipCheckAPIServerDNS = os.Getenv("KUBENURSE_CHECK_API_SERVER_DNS") == "false"
	chk.SkipCheckMeIngress = os.Getenv("KUBENURSE_CHECK_ME_INGRESS") == "false"
	chk.SkipCheckMeService = os.Getenv("KUBENURSE_CHECK_ME_SERVICE") == "false"
	chk.SkipCheckNeighbourhood = os.Getenv("KUBENURSE_CHECK_NEIGHBOURHOOD") == "false"

	chk.UseTLS = server.useTLS

	server.checker = chk

	// setup http routes
	mux.HandleFunc("/ready", server.readyHandler())
	mux.HandleFunc("/alive", server.aliveHandler())
	mux.HandleFunc("/alwayshappy", func(http.ResponseWriter, *http.Request) {})
	mux.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))
	mux.Handle("/", http.RedirectHandler("/alive", http.StatusMovedPermanently))

	return server, nil
}

// Run starts the periodic checker and the http/https server(s) and blocks until Shutdown was called.
func (s *Server) Run() error {
	var (
		wg   sync.WaitGroup
		errc = make(chan error, 2) // max two errors can happen
	)

	wg.Add(1)

	go func() {
		defer wg.Done()

		s.checker.RunScheduled(5 * time.Second)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := s.http.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				errc <- fmt.Errorf("listen http: %w", err)
			}
		}
	}()

	if s.useTLS {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := s.https.ListenAndServeTLS(
				os.Getenv("KUBENURSE_CERT_FILE"),
				os.Getenv("KUBENURSE_CERT_KEY"),
			); err != nil {
				if err != http.ErrServerClosed {
					errc <- fmt.Errorf("listen https: %w", err)
				}
			}
		}()
	}

	wg.Wait()
	close(errc)

	// return the first error if there was one
	for err := range errc {
		if err != nil {
			return err
		}
	}

	return nil
}

// Shutdown disables the readiness probe and then gracefully halts the kubenurse http/https server(s).
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.ready = false
	s.mu.Unlock()

	// stop the scheduled checker
	s.checker.StopScheduled()

	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	if s.useTLS {
		if err := s.https.Shutdown(ctx); err != nil {
			return fmt.Errorf("stop https server: %w", err)
		}
	}

	return nil
}
