package servicecheck

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptrace"
	"sync/atomic"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
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
		errorAccounted := r.Context().Value(kubenurseErrorAccountedKey{}).(*atomic.Bool)
		l := map[string]string{
			"type":  kubenurseTypeLabel,
			"event": traceEventType,
		}

		// If we get an error inside a trace, log it
		if err != nil {
			errorCounter.WithLabelValues(traceEventType, kubenurseTypeLabel).Inc()
			metrics.GetOrCreateCounter(genMetricName(errCounter, l)).Inc()
			errorAccounted.Store(true) // mark the error as accounted, so we don't increase the error counter twice.
			slog.Error("request failure in httptrace", "event_type", traceEventType, "request_type", kubenurseTypeLabel, "err", err)

			return
		}

		httpclientTraceReqDuration.WithLabelValues(traceEventType, kubenurseTypeLabel).Observe(td)
		metrics.GetOrCreatePrometheusHistogramExt(genMetricName(
			hcTraceReqDurSec, l), durHist,
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

		// typeFromCtxFn := promhttp.WithLabelFromCtx("type", func(ctx context.Context) string {
		// 	return ctx.Value(kubenurseTypeKey{}).(string)
		// })

		rt := next // variable pinning :) essential, to prevent always re-instrumenting the original variable
		// rt = promhttp.InstrumentRoundTripperCounter(httpclientReqTotal, rt, typeFromCtxFn)
		// rt = promhttp.InstrumentRoundTripperDuration(httpclientReqDuration, rt, typeFromCtxFn)

		kubenurseRequestType := r.Context().Value(kubenurseTypeKey{}).(string)
		errorAccounted := r.Context().Value(kubenurseErrorAccountedKey{}).(*atomic.Bool)

		start = time.Now()
		resp, err := rt.RoundTrip(r)
		l := addLabels("type", kubenurseRequestType, nil)

		if err == nil {
			metrics.GetOrCreateCounter(genMetricName(
				hcReqTotal, addLabels("code", fmt.Sprintf("%d", resp.StatusCode), l)),
			).Inc()
			httpclientReqTotal.With(addLabels("code", fmt.Sprintf("%d", resp.StatusCode), l)).Inc()

			metrics.GetOrCreatePrometheusHistogramExt(genMetricName(
				hcReqDurSec, l), durHist,
			).UpdateDuration(start)
			httpclientReqDuration.With(l).Observe(time.Since(start).Seconds())

			if resp.StatusCode != http.StatusOK {
				eventType := fmt.Sprintf("status_code_%d", resp.StatusCode)

				errorCounter.WithLabelValues(eventType, kubenurseRequestType).Inc()
				metrics.GetOrCreateCounter(genMetricName(errCounter, addLabels("event", eventType, l))).Inc()
				slog.Error("request failure in httptrace",
					"event_type", eventType,
					"request_type", kubenurseRequestType)
			}
		} else {
			eventType := "round_trip_error"
			metrics.GetOrCreateCounter(genMetricName(hcReqTotal, addLabels("code", eventType, l))).Inc()
			httpclientReqTotal.With(addLabels("code", eventType, l)).Inc() // also increment the total counter, as InstrumentRoundTripperCounter only instruments successful requests

			if !errorAccounted.Load() {
				errorCounter.WithLabelValues(eventType, kubenurseRequestType).Inc()
				metrics.GetOrCreateCounter(genMetricName(errCounter, addLabels("event", eventType, l))).Inc()
			}
			slog.Error("request failure in httptrace",
				"event_type", eventType,
				"request_type", kubenurseRequestType, "err", err)
		}

		return resp, err
	})
}

func addLabels(tag, value string, l map[string]string) map[string]string {
	ret := map[string]string{}
	maps.Copy(ret, l)
	ret[tag] = value
	return ret
}

func mapToLabels(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	ret := ""
	for k, v := range m {
		ret += fmt.Sprintf("%s=%q,", k, v)
	}

	return ret[:len(ret)-1]
}

func genMetricName(name string, m map[string]string) string {
	return fmt.Sprintf("%s_%s{%s}", MetricsNamespace, name, mapToLabels(m))
}
