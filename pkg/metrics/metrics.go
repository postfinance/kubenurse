package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

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
			Name:   "kubenurse_request_duration",
			Help:   "Kubenurse request duration partitioned by error type",
			MaxAge: 1 * time.Minute,
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(ErrorCounter)
	prometheus.MustRegister(DurationSummary)
}
