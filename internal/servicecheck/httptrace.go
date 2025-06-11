package servicecheck

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptrace"
	"sync/atomic"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

// unique type for context.Context to avoid collisions.
type (
	kubenurseTypeKey           struct{}
	kubenurseErrorAccountedKey struct{}
)

const (
	hcReqTotal       = "httpclient_requests_total"
	hcReqDurSec      = "httpclient_request_duration_seconds"
	hcTraceReqDurSec = "httpclient_trace_request_duration_seconds"
	errCounter       = "errors_total"
)

// RoundTripperFunc is a function which performs a round-trip check and potentially returns a response/error
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

// This collects traces and logs errors. As promhttp.InstrumentRoundTripperTrace doesn't process
// errors, this is custom made and inspired by prometheus/client_golang's promhttp
//
//nolint:funlen // needed to pack all histograms and use them directly in the httptrace wrapper
func withHttptrace(next http.RoundTripper, durHist []float64) http.RoundTripper {
	collectMetric := func(traceEventType string, start time.Time, r *http.Request, err error) {
		kubenurseTypeLabel := r.Context().Value(kubenurseTypeKey{}).(string)
		errorAccounted := r.Context().Value(kubenurseErrorAccountedKey{}).(*atomic.Bool)
		l := []string{"type", kubenurseTypeLabel, "event", traceEventType}

		// If we get an error inside a trace, log it
		if err != nil {
			metrics.GetOrCreateCounter(genMetricName(errCounter, l...)).Inc()
			errorAccounted.Store(true) // mark the error as accounted, so we don't increase the error counter twice.
			slog.Error("request failure in httptrace", "event_type", traceEventType, "request_type", kubenurseTypeLabel, "err", err)

			return
		}

		metrics.GetOrCreatePrometheusHistogramExt(genMetricName(
			hcTraceReqDurSec, l...), durHist,
		).UpdateDuration(start)
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

		rt := next // variable pinning :) essential, to prevent always re-instrumenting the original variable

		kubenurseRequestType := r.Context().Value(kubenurseTypeKey{}).(string)
		errorAccounted := r.Context().Value(kubenurseErrorAccountedKey{}).(*atomic.Bool)

		start = time.Now()
		resp, err := rt.RoundTrip(r)
		l := []string{"type", kubenurseRequestType}

		if err == nil {
			metrics.GetOrCreateCounter(genMetricName(
				hcReqTotal, append(l, "code", fmt.Sprintf("%d", resp.StatusCode))...),
			).Inc()

			metrics.GetOrCreatePrometheusHistogramExt(genMetricName(
				hcReqDurSec, l...), durHist,
			).UpdateDuration(start)

			if resp.StatusCode != http.StatusOK {
				eventType := fmt.Sprintf("status_code_%d", resp.StatusCode)

				metrics.GetOrCreateCounter(genMetricName(errCounter, append(l, "event", eventType)...)).Inc()
				slog.Error("request failure in httptrace",
					"event_type", eventType,
					"request_type", kubenurseRequestType)
			}
		} else {
			eventType := "round_trip_error"
			metrics.GetOrCreateCounter(genMetricName(hcReqTotal, append(l, "code", eventType)...)).Inc()

			if !errorAccounted.Load() {
				metrics.GetOrCreateCounter(genMetricName(errCounter, append(l, "event", eventType)...)).Inc()
			}
			slog.Error("request failure in httptrace",
				"event_type", eventType,
				"request_type", kubenurseRequestType, "err", err)
		}

		return resp, err
	})
}

func genMetricName(name string, kvs ...string) string {
	n := len(kvs)
	labels := ""
	if n > 0 {
		if n%2 != 0 {
			panic("odd number or label tags, cannot construct the metric name")
		}
		for i := 0; i < n; i += 2 {
			labels += fmt.Sprintf("%s=%q,", kvs[i], kvs[i+1])
		}
		labels = labels[:len(labels)-1]
	}

	return fmt.Sprintf("%s_%s{%s}", MetricsNamespace, name, labels)
}
