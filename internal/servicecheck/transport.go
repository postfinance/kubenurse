package servicecheck

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	//nolint:gosec // This is the well-known path to Kubernetes serviceaccount tokens.
	tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	caFile    = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

// doRequest does an http request only to get the http status code
func (c *Checker) doRequest(url string) (string, error) {
	// Read Bearer Token file from ServiceAccount
	token, err := os.ReadFile(tokenFile)
	if err != nil {
		return errStr, fmt.Errorf("load kubernetes serviceaccount token from %s: %w", tokenFile, err)
	}

	req, _ := http.NewRequest("GET", url, http.NoBody)

	// Only add the Bearer for API Server Requests
	if strings.HasSuffix(url, "/version") {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err.Error(), err
	}

	// Body is non-nil if err is nil, so close it
	_ = resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return okStr, nil
	}

	return resp.Status, errors.New(resp.Status)
}

// generateRoundTripper returns a custom http.RoundTripper, including the k8s CA.
func generateRoundTripper(extraCA string, insecure bool) (http.RoundTripper, error) {
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
	if extraCA != "" {
		caCert, err := os.ReadFile(extraCA) // Intentionally included by the user.
		if err != nil {
			return nil, fmt.Errorf("could not load certificate %s: %w", extraCA, err)
		}

		if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
			return nil, errors.New("could not append extra ca cert to system certpool")
		}
	}

	// Configure transport
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecure, //nolint:gosec // Can be true if the user requested this.
		RootCAs:            rootCAs,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return transport, nil
}
