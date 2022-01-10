package servicecheck

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func withRequestTracing(registry *prometheus.Registry, transport http.RoundTripper) http.RoundTripper {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "httpclient_requests_total",
			Help:      "A counter for requests from the kubenurse http client.",
		},
		[]string{"code", "method"},
	)

	latencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "httpclient_trace_request_duration_seconds",
			Help:      "Latency histogram for requests from the kubenurse http client. Time in seconds since the start of the http request.",
			Buckets:   []float64{.0005, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"event"},
	)

	// histVec has no labels, making it a zero-dimensional ObserverVec.
	histVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "httpclient_request_duration_seconds",
			Help:      "A latency histogram of request latencies from the kubenurse http client.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{},
	)

	// Register all of the metrics in the standard registry.
	registry.MustRegister(counter, latencyVec, histVec)

	// Define functions for the available httptrace.ClientTrace hook
	// functions that we want to instrument.
	trace := &promhttp.InstrumentTrace{
		DNSStart: func(t float64) {
			latencyVec.WithLabelValues("dns_start").Observe(t)
		},
		DNSDone: func(t float64) {
			latencyVec.WithLabelValues("dns_done").Observe(t)
		},
		ConnectStart: func(t float64) {
			latencyVec.WithLabelValues("connect_start").Observe(t)
		},
		ConnectDone: func(t float64) {
			latencyVec.WithLabelValues("connect_done").Observe(t)
		},
		TLSHandshakeStart: func(t float64) {
			latencyVec.WithLabelValues("tls_handshake_start").Observe(t)
		},
		TLSHandshakeDone: func(t float64) {
			latencyVec.WithLabelValues("tls_handshake_done").Observe(t)
		},
		WroteRequest: func(t float64) {
			latencyVec.WithLabelValues("wrote_request").Observe(t)
		},
		GotFirstResponseByte: func(t float64) {
			latencyVec.WithLabelValues("got_first_resp_byte").Observe(t)
		},
	}

	// Wrap the default RoundTripper with middleware.
	roundTripper := promhttp.InstrumentRoundTripperCounter(counter,
		promhttp.InstrumentRoundTripperTrace(trace,
			promhttp.InstrumentRoundTripperDuration(histVec,
				transport,
			),
		),
	)

	return roundTripper
}
