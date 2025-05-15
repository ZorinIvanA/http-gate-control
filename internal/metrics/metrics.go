package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gate_requests_total",
			Help: "Total number of gate requests.",
		},
	)

	// Меняем Counter на Gauge для возможности сброса
	BlockedCounter = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gate_blocked_total",
			Help: "Number of blocked gate openings in the current window.",
		},
	)

	RateLimitCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of rate limit exceeded errors.",
		},
	)
)

func init() {
	prometheus.MustRegister(RequestCounter, BlockedCounter, RateLimitCounter)
}

// MustRegister — регистрирует метрики вручную, если нужно
func MustRegister() {
	prometheus.MustRegister(RequestCounter, BlockedCounter, RateLimitCounter)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RequestCounter.Inc()
		next.ServeHTTP(w, r)
	})
}
