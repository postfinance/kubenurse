// Package kubenurse contains the server code for the kubenurse service.
package kubenurse

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/postfinance/kubenurse/internal/servicecheck"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const defaultCheckInterval = 5 * time.Second

// Server is used to build the kubenurse http/https server(s).
type Server struct {
	http  http.Server
	https http.Server

	checker *servicecheck.Checker

	// Configuration options
	useTLS        bool
	checkInterval time.Duration
	// If we want to consider kubenurses on unschedulable nodes
	allowUnschedulable bool

	ready atomic.Bool

	neighboursTTLCache TTLCache[string]
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
// * KUBENURSE_NEIGHBOUR_LIMIT
// * KUBENURSE_SHUTDOWN_DURATION
// * KUBENURSE_CHECK_API_SERVER_DIRECT
// * KUBENURSE_CHECK_API_SERVER_DNS
// * KUBENURSE_CHECK_ME_INGRESS
// * KUBENURSE_CHECK_ME_SERVICE
// * KUBENURSE_CHECK_NEIGHBOURHOOD
// * KUBENURSE_CHECK_INTERVAL
func New(c client.Client) (*Server, error) { //nolint:funlen // TODO: use a flag parsing library (e.g. ff) to reduce complexity
	mux := http.NewServeMux()

	checkInterval := defaultCheckInterval

	if v, ok := os.LookupEnv("KUBENURSE_CHECK_INTERVAL"); ok {
		var err error
		checkInterval, err = time.ParseDuration(v)
		if err != nil {
			return nil, err
		}
	}

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
		checkInterval:      checkInterval,
		ready:              atomic.Bool{},
	}

	server.ready.Store(true)
	server.neighboursTTLCache.Init(60 * time.Second)

	var histogramBuckets []float64

	if bucketsString := os.Getenv("KUBENURSE_HISTOGRAM_BUCKETS"); bucketsString != "" {
		for bucketStr := range strings.SplitSeq(bucketsString, ",") {
			bucket, err := strconv.ParseFloat(bucketStr, 64)
			if err != nil {
				slog.Error("couldn't parse one of the custom histogram buckets", "bucket", bucket, "err", err)
				os.Exit(1)
			}

			histogramBuckets = append(histogramBuckets, bucket)
		}

		err := metrics.ValidateBuckets(histogramBuckets)
		if err != nil {
			slog.Error("custom histogram buckets validation failed", "bucket_bounds", histogramBuckets, "err", err)
			os.Exit(1)
		}
	}

	// setup checker
	chk, err := servicecheck.New(c, server.allowUnschedulable, 1*time.Second, func(s string) servicecheck.Histogram {
		if os.Getenv("KUBENURSE_VICTORIAMETRICS_HISTOGRAM") == "true" {
			return metrics.GetOrCreateHistogram(s)
		} else {
			if histogramBuckets == nil {
				histogramBuckets = metrics.PrometheusHistogramDefaultBuckets
			}
			return metrics.GetOrCreatePrometheusHistogramExt(s, histogramBuckets)
		}
	})
	if err != nil {
		return nil, err
	}

	shutdownDuration := 5 * time.Second

	if v, ok := os.LookupEnv("KUBENURSE_SHUTDOWN_DURATION"); ok {
		shutdownDuration, err = time.ParseDuration(v)
		if err != nil {
			return nil, err
		}
	}

	chk.ShutdownDuration = shutdownDuration
	chk.KubenurseIngressURL = os.Getenv("KUBENURSE_INGRESS_URL")
	chk.KubenurseServiceURL = os.Getenv("KUBENURSE_SERVICE_URL")
	chk.KubernetesServiceHost = os.Getenv("KUBERNETES_SERVICE_HOST")
	chk.KubernetesServicePort = os.Getenv("KUBERNETES_SERVICE_PORT")
	chk.KubernetesServiceDNS = getOrDefault("KUBERNETES_SERVICE_DNS", "kubernetes.default.svc.cluster.local")
	chk.KubenurseNamespace = os.Getenv("KUBENURSE_NAMESPACE")
	chk.NeighbourFilter = os.Getenv("KUBENURSE_NEIGHBOUR_FILTER")
	neighLimit := os.Getenv("KUBENURSE_NEIGHBOUR_LIMIT")

	if neighLimit != "" {
		chk.NeighbourLimit, err = strconv.Atoi(neighLimit)
		if err != nil {
			return nil, err
		}
	} else {
		chk.NeighbourLimit = 10
	}

	//nolint:goconst // No need to make "false" a constant in my opinion, readability is better like this.
	chk.SkipCheckAPIServerDirect = os.Getenv("KUBENURSE_CHECK_API_SERVER_DIRECT") == "false"
	chk.SkipCheckAPIServerDNS = os.Getenv("KUBENURSE_CHECK_API_SERVER_DNS") == "false"
	chk.SkipCheckMeIngress = os.Getenv("KUBENURSE_CHECK_ME_INGRESS") == "false"
	chk.SkipCheckMeService = os.Getenv("KUBENURSE_CHECK_ME_SERVICE") == "false"
	chk.SkipCheckNeighbourhood = os.Getenv("KUBENURSE_CHECK_NEIGHBOURHOOD") == "false"

	chk.UseTLS = server.useTLS

	// Extra checks parsing
	if extraChecks := os.Getenv("KUBENURSE_EXTRA_CHECKS"); extraChecks != "" {
		for _, extraCheck := range strings.Split(extraChecks, "|") {
			requestType, url, fnd := strings.Cut(extraCheck, ":")
			if !fnd {
				slog.Error("couldn't parse one of extraChecks", "extraCheck", extraCheck)
				return nil, fmt.Errorf("extra checks parsing - missing colon ':' between metric name and url")
			}

			chk.ExtraChecks[requestType] = url
		}
	}

	server.checker = chk

	// setup http routes
	mux.HandleFunc("/ready", server.readyHandler())
	mux.HandleFunc("/alive", server.aliveHandler())
	mux.HandleFunc("/alwayshappy", server.alwaysHappyHandler())
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	mux.Handle("/", http.RedirectHandler("/alive", http.StatusMovedPermanently))

	return server, nil
}

// Run starts the periodic checker and the http/https server(s) and blocks until Shutdown was called.
func (s *Server) Run(ctx context.Context) error {
	var (
		wg   sync.WaitGroup
		errc = make(chan error, 2) // max two errors can happen
	)

	go func() { // update the incoming neighbouring check gauge every second
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()

		for range t.C {
			metrics.GetOrCreateGauge("kubenurse_neighbourhood_incoming_checks", nil).Set(
				float64(s.neighboursTTLCache.ActiveEntries()),
			)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(s.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.checker.Run(ctx)
			case <-ctx.Done():
				return
			}
		}
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

	slog.Info("kubenurse just started")

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
func (s *Server) Shutdown() error {
	s.ready.Store(false)

	// wait before actually shutting down the http/s server, as the updated
	// endpoints for the kubenurse service might not have propagated everywhere
	// (other kubenurse/ingress controller) yet, which will lead to
	// me_ingress or path errors in other pods
	time.Sleep(s.checker.ShutdownDuration)

	// background ctx since, the "root" context is already canceled
	ctx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

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

func getOrDefault(envVar, defaultVal string) string {
	if val := os.Getenv(envVar); val != "" {
		return val
	}

	return defaultVal
}
