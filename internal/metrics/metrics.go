package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	rpsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_gate_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method"},
	)
	blockedCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "http_gate_blocked_total",
			Help: "Total number of blocked gate openings.",
		},
	)
)

func NewMetrics() Metrics {
	return &prometheusMetrics{}
}

type Metrics interface {
	IncRPS(method string)
	IncBlocked()
	Handler() http.Handler
}

type prometheusMetrics struct{}

func (m *prometheusMetrics) IncRPS(method string) {
	rpsCounter.WithLabelValues(method).Inc()
}

func (m *prometheusMetrics) IncBlocked() {
	blockedCounter.Inc()
}

func (m *prometheusMetrics) Handler() http.Handler {
	return promhttp.Handler()
}
