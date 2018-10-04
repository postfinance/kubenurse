package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/postfinance/kubenurse/pkg/checker"
	"github.com/postfinance/kubenurse/pkg/kubediscovery"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	chk = &checker.Checker{}
)

const (
	caFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	nurse  = "I'm ready to help you!"
)

func main() {
	// Setup http transport
	transport, err := GenerateRoundTripper()
	if err != nil {
		log.Printf("using default transport: %s", err)
		transport = http.DefaultTransport
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	// Setup checker
	chk.KubenurseIngressUrl = os.Getenv("KUBENURSE_SERVICE_URL")
	chk.KubenurseServiceUrl = os.Getenv("KUBENURSE_SERVICE_URL")
	chk.KubernetesServiceHost = os.Getenv("KUBERNETES_SERVICE_HOST")
	chk.KubernetesServicePort = os.Getenv("KUBERNETES_SERVICE_PORT")
	chk.KubeNamespace = os.Getenv("KUBE_NAMESPACE")
	chk.NeighbourFilter = os.Getenv("KUBENURSE_NEIGHBOUR_FILTER")
	chk.HttpClient = client

	// Setup http routes
	http.HandleFunc("/alive", aliveHandler)
	http.HandleFunc("/alwayshappy", func(http.ResponseWriter, *http.Request) {})
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", http.RedirectHandler("/alive", http.StatusMovedPermanently))

	fmt.Println(nurse) // most important line of this project

	// Start listener and checker
	go func() {
		chk.RunScheduled(5 * time.Second)
		log.Fatalln("checker exited")
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func aliveHandler(w http.ResponseWriter, r *http.Request) {
	type Output struct {
		Hostname string              `json:"hostname"`
		Headers  map[string][]string `json:"headers"`

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
		Neighbourhood:      res.Neighbourhood,
		NeighbourhoodState: res.NeighbourhoodState,
	}
	out.Hostname, _ = os.Hostname()

	// Generate output output
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	enc.Encode(out)

}

func GenerateRoundTripper() (http.RoundTripper, error) {
	insecureEnv := os.Getenv("KUBENURSE_INSECURE")
	insecure, _ := strconv.ParseBool(insecureEnv)

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("could not load certificate %s: %s", caFile, err)
	}

	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("could not append ca cert to system certpool")
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecure,
		RootCAs:            rootCAs,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return transport, nil
}
