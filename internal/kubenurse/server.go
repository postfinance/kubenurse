// Package kubenurse contains the server code for the kubenurse service.
package kubenurse

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/postfinance/kubenurse/pkg/checker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server is used to build the kubenurse http/https server(s).
type Server struct {
	http  http.Server
	https http.Server

	checker *checker.Checker

	// Configuration options
	useTLS bool
	// If we want to consider kubenurses on unschedulable nodes
	allowUnschedulable bool
	extraCA            string
	insecure           bool

	// Mutex to protect ready flag
	mu    *sync.Mutex
	ready bool
}

// New creates a new kubenurse server. The server can be configured with the following environment variables:
// * KUBENURSE_USE_TLS
// * KUBENURSE_ALLOW_UNSCHEDULABL
// * KUBENURSE_INGRESS_URL
// * KUBENURSE_SERVICE_URL
// * KUBERNETES_SERVICE_HOST
// * KUBERNETES_SERVICE_PORT
// * KUBENURSE_NAMESPACE
// * KUBENURSE_NEIGHBOUR_FILTER
// * KUBENURSE_EXTRA_CA
// * KUBENURSE_INSECURE
func New() (*Server, error) {
	mux := http.NewServeMux()

	server := &Server{
		http: http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		https: http.Server{
			Addr:    ":8443",
			Handler: mux,
		},

		//nolint:goconst // No need to make "true" a constant in my opinion, readability is better like this.
		useTLS:             os.Getenv("KUBENURSE_USE_TLS") == "true",
		allowUnschedulable: os.Getenv("KUBENURSE_ALLOW_UNSCHEDULABLE") == "true",
		extraCA:            os.Getenv("KUBENURSE_EXTRA_CA"),
		insecure:           os.Getenv("KUBENURSE_INSECURE") == "true",

		mu:    new(sync.Mutex),
		ready: true,
	}

	// setup http transport
	transport, err := server.generateRoundTripper()
	if err != nil {
		log.Printf("using default transport: %s", err)

		transport = http.DefaultTransport
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	// setup checker
	chk, err := checker.New(context.TODO(), client, 3*time.Second, server.allowUnschedulable)
	if err != nil {
		return nil, err
	}

	chk.KubenurseIngressURL = os.Getenv("KUBENURSE_INGRESS_URL")
	chk.KubenurseServiceURL = os.Getenv("KUBENURSE_SERVICE_URL")
	chk.KubernetesServiceHost = os.Getenv("KUBERNETES_SERVICE_HOST")
	chk.KubernetesServicePort = os.Getenv("KUBERNETES_SERVICE_PORT")
	chk.KubenurseNamespace = os.Getenv("KUBENURSE_NAMESPACE")
	chk.NeighbourFilter = os.Getenv("KUBENURSE_NEIGHBOUR_FILTER")
	chk.UseTLS = server.useTLS

	server.checker = chk

	// setup http routes
	mux.HandleFunc("/ready", server.readyHandler())
	mux.HandleFunc("/alive", server.aliveHandler())
	mux.HandleFunc("/alwayshappy", func(http.ResponseWriter, *http.Request) {})
	mux.Handle("/metrics", promhttp.Handler())
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
		log.Printf("checker exited")
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
