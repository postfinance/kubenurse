package servicecheck

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// TODO:
// - RoundTripperCounter and RoundTripper duration useful? Was never officially documented and I don't see anything usable with it

// unique type for context.Context to avoid collisions.
type kubenurseContextKey struct{}

//http.RoundTripper
// TODO: Easier method to get a round tripper?
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

//
func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

// Ensure RoundTripperFunc is a http.RoundTripper
var _ http.RoundTripper = (*RoundTripperFunc)(nil)

// TODO: Description
// This collects traces and logs errors. As promhttp.InstrumentRoundTripperTrace doesn't process
// errors, this is custom made and inspired by prometheus/client_golang's promhttp
func withHttptrace(registry *prometheus.Registry, next http.RoundTripper, latencyVec *prometheus.HistogramVec) http.RoundTripper {
	collectMetric := func(traceType string, start time.Time, r *http.Request, err error) {
		td := time.Since(start).Seconds()
		kubenurseCheckLabel := r.Context().Value(kubenurseContextKey{}).(string)

		// If we got an error inside a trace, log it and do not collect metrics
		if err != nil {
			log.Printf("httptrace: failed %s for %s with %v", traceType, kubenurseCheckLabel, err)
			return
		}

		latencyVec.WithLabelValues(traceType, kubenurseCheckLabel).Observe(td)
	}

	// Return a http.RoundTripper for tracing requests
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		// Capture request time
		start := time.Now()

		// Add tracing hooks
		trace := &httptrace.ClientTrace{
			GotConn: func(info httptrace.GotConnInfo) {
				collectMetric("got_conn", start, r, nil)
			},
			DNSStart: func(info httptrace.DNSStartInfo) {
				collectMetric("dns_start", start, r, nil)
			},
			DNSDone: func(info httptrace.DNSDoneInfo) {
				collectMetric("dns_done", start, r, info.Err)
			},
			ConnectStart: func(_, _ string) {
				collectMetric("connect_start", start, r, nil)
			},
			ConnectDone: func(_, _ string, err error) {
				collectMetric("connect_done", start, r, err)
			},
			TLSHandshakeStart: func() {
				collectMetric("tls_handshake_start", start, r, nil)
			},
			TLSHandshakeDone: func(_ tls.ConnectionState, err error) {
				collectMetric("tls_handshake_done", start, r, nil)
			},
			WroteRequest: func(info httptrace.WroteRequestInfo) {
				collectMetric("wrote_request", start, r, info.Err)
			},
			GotFirstResponseByte: func() {
				collectMetric("got_first_resp_byte", start, r, nil)
			},
		}

		// Do request with tracing enabled
		r = r.WithContext(httptrace.WithClientTrace(r.Context(), trace))

		return next.RoundTrip(r)
	})
}
