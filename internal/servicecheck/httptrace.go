package servicecheck

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// unique type for context.Context to avoid collisions.
type kubenurseTypeKey struct{}

// RoundTripperFunc is a function which performs a round-trip check and potentially returns a response/error
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

// This collects traces and logs errors. As promhttp.InstrumentRoundTripperTrace doesn't process
// errors, this is custom made and inspired by prometheus/client_golang's promhttp
//
//nolint:funlen // needed to pack all histograms and use them directly in the httptrace wrapper
func withHttptrace(registry *prometheus.Registry, next http.RoundTripper, durHist []float64) http.RoundTripper {
	httpclientReqTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Name:      "httpclient_requests_total",
			Help:      "A counter for requests from the kubenurse http client.",
		},
		[]string{"code", "type"},
	)

	httpclientReqDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Name:      "httpclient_request_duration_seconds",
			Help:      "A latency histogram of request latencies from the kubenurse http client.",
			Buckets:   durHist,
		},
		[]string{"type"},
	)

	httpclientTraceReqDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Name:      "httpclient_trace_request_duration_seconds",
			Help:      "Latency histogram for requests from the kubenurse http client. Time in seconds since the start of the http request.",
			Buckets:   durHist,
		},
		[]string{"event", "type"},
	)

	errorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Name:      "errors_total",
			Help:      "Kubenurse error counter partitioned by error type",
		},
		[]string{"event", "type"},
	)

	registry.MustRegister(httpclientReqTotal, httpclientReqDuration, httpclientTraceReqDuration, errorCounter)

	collectMetric := func(traceEventType string, start time.Time, r *http.Request, err error) {
		td := time.Since(start).Seconds()
		kubenurseTypeLabel := r.Context().Value(kubenurseTypeKey{}).(string)

		// If we get an error inside a trace, log it
		if err != nil {
			errorCounter.WithLabelValues(traceEventType, kubenurseTypeLabel).Inc()
			slog.Error("request failure in httptrace", "event_type", traceEventType, "request_type", kubenurseTypeLabel, "err", err)

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
			TLSHandshakeDone: func(_ tls.ConnectionState, err error) {
				collectMetric("tls_handshake_done", start, r, err)
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

		kubenurseRequestType := r.Context().Value(kubenurseTypeKey{}).(string)
		resp, err := rt.RoundTrip(r)

		if err == nil {
			if resp.StatusCode != http.StatusOK {
				eventType := fmt.Sprintf("status_code_%d", resp.StatusCode)

				errorCounter.WithLabelValues(eventType, kubenurseRequestType).Inc()
				slog.Error("request failure in httptrace",
					"event_type", eventType,
					"request_type", kubenurseRequestType)
			}
		} else {
			eventType := "round_trip_error"
			labels := map[string]string{
				"code": eventType, // we reuse round_trip_error as status code to prevent introducing a new label
				"type": kubenurseRequestType,
			}
			httpclientReqTotal.With(labels).Inc() // also increment the total counter, as InstrumentRoundTripperCounter only instruments successful requests
			// errorCounter.WithLabelValues(eventType, kubenurseRequestType).Inc()
			// normally, errors are already accounted for in the ClientTrace section.
			// we still log the error, so in the future we can compare the log entries and see if somehow
			// an error isn't catched in the ClientTrace section
			slog.Error("request failure in httptrace",
				"event_type", eventType,
				"request_type", kubenurseRequestType, "err", err)
		}

		return resp, err
	})
}
