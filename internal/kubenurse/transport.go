package kubenurse

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
)

const (
	caFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

// generateRoundTripper returns a custom http.RoundTripper, including the k8s
// CA. If env KUBENURSE_INSECURE is set to true, certificates are not validated.
func (s *Server) generateRoundTripper() (http.RoundTripper, error) {
	// Append default certpool
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Append ServiceAccount cacert
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("could not load certificate %s: %w", caFile, err)
	}

	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("could not append ca cert to system certpool")
	}

	// Append extra CA, if set
	if s.extraCA != "" {
		caCert, err := os.ReadFile(s.extraCA)
		if err != nil {
			return nil, fmt.Errorf("could not load certificate %s: %w", s.extraCA, err)
		}

		if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
			return nil, errors.New("could not append extra ca cert to system certpool")
		}
	}

	// Configure transport
	tlsConfig := &tls.Config{
		InsecureSkipVerify: s.insecure, //nolint:gosec // Can be true if the user requested this.
		RootCAs:            rootCAs,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return transport, nil
}
