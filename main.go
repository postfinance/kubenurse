package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/postfinance/kubenurse/pkg/checker"
	"github.com/postfinance/kubenurse/pkg/kubediscovery"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	caFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	nurse  = "I'm ready to help you!"
)

//nolint:funlen
func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	serverTLS := http.Server{
		Addr:    ":8443",
		Handler: mux,
	}
	useTLS := os.Getenv("KUBENURSE_USE_TLS") == "true"

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case s := <-sig:
			log.Printf("shutting down, received signal %s", s)

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer shutdownCancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Fatalln(err)
			}

			if useTLS {
				if err := serverTLS.Shutdown(shutdownCtx); err != nil {
					log.Fatalln(err)
				}
			}

			cancel()
		case <-ctx.Done():
		}
	}()

	// setup http transport
	transport, err := GenerateRoundTripper()
	if err != nil {
		log.Printf("using default transport: %s", err)

		transport = http.DefaultTransport
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	// If we want to consider kubenurses on unschedulable nodes
	allowUnschedulable := os.Getenv("KUBENURSE_ALLOW_UNSCHEDULABLE") == "true"

	// setup checker
	chk, err := checker.New(ctx, client, 3*time.Second, allowUnschedulable)
	if err != nil {
		log.Fatalln(err)
	}

	chk.KubenurseIngressURL = os.Getenv("KUBENURSE_INGRESS_URL")
	chk.KubenurseServiceURL = os.Getenv("KUBENURSE_SERVICE_URL")
	chk.KubernetesServiceHost = os.Getenv("KUBERNETES_SERVICE_HOST")
	chk.KubernetesServicePort = os.Getenv("KUBERNETES_SERVICE_PORT")
	chk.KubenurseNamespace = os.Getenv("KUBENURSE_NAMESPACE")
	chk.NeighbourFilter = os.Getenv("KUBENURSE_NEIGHBOUR_FILTER")
	chk.UseTLS = useTLS

	// setup http routes
	mux.HandleFunc("/alive", aliveHandler(chk))
	mux.HandleFunc("/alwayshappy", func(http.ResponseWriter, *http.Request) {})
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", http.RedirectHandler("/alive", http.StatusMovedPermanently))

	fmt.Println(nurse) // most important line of this project

	// Start listener and checker
	go func() {
		chk.RunScheduled(5 * time.Second)
		log.Fatalln("checker exited")
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalln(err)
			}
		}
	}()

	if useTLS {
		go func() {
			if err := serverTLS.ListenAndServeTLS(os.Getenv("KUBENURSE_CERT_FILE"), os.Getenv("KUBENURSE_CERT_KEY")); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalln(err)
				}
			}
		}()
	}

	<-ctx.Done()
}

func aliveHandler(chk *checker.Checker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type Output struct {
			Hostname   string              `json:"hostname"`
			Headers    map[string][]string `json:"headers"`
			UserAgent  string              `json:"user_agent"`
			RequestURI string              `json:"request_uri"`
			RemoteAddr string              `json:"remote_addr"`

			// checker.Result
			APIServerDirect string `json:"api_server_direct"`
			APIServerDNS    string `json:"api_server_dns"`
			MeIngress       string `json:"me_ingress"`
			MeService       string `json:"me_service"`

			// kubediscovery
			NeighbourhoodState string                    `json:"neighbourhood_state"`
			Neighbourhood      []kubediscovery.Neighbour `json:"neighbourhood"`
		}

		// Run checks now
		res, haserr := chk.Run()
		if haserr {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Add additional data
		out := Output{
			APIServerDNS:       res.APIServerDNS,
			APIServerDirect:    res.APIServerDirect,
			MeIngress:          res.MeIngress,
			MeService:          res.MeService,
			Headers:            r.Header,
			UserAgent:          r.UserAgent(),
			RequestURI:         r.RequestURI,
			RemoteAddr:         r.RemoteAddr,
			Neighbourhood:      res.Neighbourhood,
			NeighbourhoodState: res.NeighbourhoodState,
		}
		out.Hostname, _ = os.Hostname()

		// Generate output output
		enc := json.NewEncoder(w)
		enc.SetIndent("", " ")
		_ = enc.Encode(out)
	}
}

// GenerateRoundTripper returns a custom http.RoundTripper, including the k8s
// CA. If env KUBENURSE_INSECURE is set to true, certificates are not validated.
func GenerateRoundTripper() (http.RoundTripper, error) {
	// Parse environment variables
	extraCA := os.Getenv("KUBENURSE_EXTRA_CA")
	insecureEnv := os.Getenv("KUBENURSE_INSECURE")
	insecure, _ := strconv.ParseBool(insecureEnv)

	// Append default certpool
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Append ServiceAccount cacert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("could not load certificate %s: %s", caFile, err)
	}

	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("could not append ca cert to system certpool")
	}

	// Append extra CA, if set
	if extraCA != "" {
		//nolint:gosec
		caCert, err := ioutil.ReadFile(extraCA)

		if err != nil {
			return nil, fmt.Errorf("could not load certificate %s: %s", extraCA, err)
		}

		if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
			return nil, errors.New("could not append extra ca cert to system certpool")
		}
	}

	// Configure transport
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecure, //nolint:gosec
		RootCAs:            rootCAs,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return transport, nil
}
