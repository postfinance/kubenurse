package servicecheck

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
)

const (
	//nolint:gosec // This is the well-known path to Kubernetes serviceaccount tokens.
	K8sTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	k8sCAFile    = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

// doRequest does an http request only to get the http status code
func (c *Checker) doRequest(ctx context.Context, url string, addOriginHeader bool) string {
	// Read Bearer Token file from ServiceAccount
	token, err := os.ReadFile(K8sTokenFile)
	if !testing.Testing() && err != nil {
		slog.Error("error in doRequest while reading k8sTokenFile", "err", err)
		return errStr
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)

	// Only add the Bearer for API Server Requests
	if strings.HasSuffix(url, "/version") {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	if addOriginHeader {
		hostname, _ := os.Hostname()
		req.Header.Add(NeighbourOriginHeader, hostname)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err.Error()
	}

	// Body is non-nil if err is nil, so close it
	_ = resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return okStr
	}

	return resp.Status
}

// generateTLSConfig returns a TLSConfig including K8s CA and the user-defined extraCA
func generateTLSConfig(extraCA string) (*tls.Config, error) {
	// Append default certpool
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Append ServiceAccount cacert
	caCert, err := os.ReadFile(k8sCAFile)
	if err != nil {
		return nil, fmt.Errorf("could not load certificate %s: %w", k8sCAFile, err)
	}

	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("could not append ca cert to system certpool")
	}

	// Append extra CA, if set
	if extraCA != "" {
		caCert, err := os.ReadFile(extraCA) //nolint:gosec // Intentionally included by the user.
		if err != nil {
			return nil, fmt.Errorf("could not load certificate %s: %w", extraCA, err)
		}

		if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
			return nil, errors.New("could not append extra ca cert to system certpool")
		}
	}

	// Configure transport
	tlsConfig := &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
	}

	return tlsConfig, nil
}
