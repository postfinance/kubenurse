// Package metrics sets-up the metrics which will be exported by kubenurse. TODO: rewrite this package.
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//nolint:gochecknoglobals
var (
	// ErrorCounter provides the kubenurse_errors_total metric
	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubenurse_errors_total",
			Help: "Kubenurse error counter partitioned by error type",
		},
		[]string{"type"},
	)

	// DurationSummary provides the kubenurse_request_duration metric
	DurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "kubenurse_request_duration",
			Help:       "Kubenurse request duration partitioned by error type",
			MaxAge:     1 * time.Minute,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"type"},
	)
)

//nolint:gochecknoinits
func init() {
	prometheus.MustRegister(ErrorCounter)
	prometheus.MustRegister(DurationSummary)
}
