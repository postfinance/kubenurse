package servicecheck

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// unique type for context.Context to avoid collisions.
type kubenurseTypeKey struct{}

// // http.RoundTripper
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

// This collects traces and logs errors. As promhttp.InstrumentRoundTripperTrace doesn't process
// errors, this is custom made and inspired by prometheus/client_golang's promhttp
func withHttptrace(registry *prometheus.Registry, next http.RoundTripper, durationHistogram []float64) http.RoundTripper {
	httpclientReqTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Name:      "httpclient_requests_total",
			Help:      "A counter for requests from the kubenurse http client.",
		},
		[]string{"code", "method", "type"},
	)

	httpclientReqDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Name:      "httpclient_request_duration_seconds",
			Help:      "A latency histogram of request latencies from the kubenurse http client.",
			Buckets:   durationHistogram,
		},
		[]string{"type"},
	)

	httpclientTraceReqDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Name:      "httpclient_trace_request_duration_seconds",
			Help:      "Latency histogram for requests from the kubenurse http client. Time in seconds since the start of the http request.",
			Buckets:   durationHistogram,
		},
		[]string{"event", "type"},
	)

	registry.MustRegister(httpclientReqTotal, httpclientReqDuration, httpclientTraceReqDuration)

	collectMetric := func(traceEventType string, start time.Time, r *http.Request, err error) {
		td := time.Since(start).Seconds()
		kubenurseTypeLabel := r.Context().Value(kubenurseTypeKey{}).(string)

		// If we got an error inside a trace, log it and do not collect metrics
		if err != nil {
			log.Printf("httptrace: failed %s for %s with %v", traceEventType, kubenurseTypeLabel, err)
			return
		}

		httpclientTraceReqDuration.WithLabelValues(traceEventType, kubenurseTypeLabel).Observe(td)
	}

	// Return a http.RoundTripper for tracing requests
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		// Capture request time
		start := time.Now()

		// Add tracing hooks
		trace := &httptrace.ClientTrace{
			GotConn: func(_ httptrace.GotConnInfo) {
				collectMetric("got_conn", start, r, nil)
			},
			DNSStart: func(_ httptrace.DNSStartInfo) {
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
			TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
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

		typeFromCtxFn := promhttp.WithLabelFromCtx("type", func(ctx context.Context) string {
			return ctx.Value(kubenurseTypeKey{}).(string)
		})

		rt := next // variable pinning :) essential, to prevent always re-instrumenting the original variable
		rt = promhttp.InstrumentRoundTripperCounter(httpclientReqTotal, rt, typeFromCtxFn)
		rt = promhttp.InstrumentRoundTripperDuration(httpclientReqDuration, rt, typeFromCtxFn)

		return rt.RoundTrip(r)
	})
}
