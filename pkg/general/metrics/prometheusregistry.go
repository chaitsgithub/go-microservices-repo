package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type prometheusRegistry struct {
	registerer *prometheus.Registerer
}

func NewPrometheusRegistry() *prometheusRegistry {
	return &prometheusRegistry{registerer: &prometheus.DefaultRegisterer}
}

func (r *prometheusRegistry) Counter(name, help string, labels ...string) Counter {
	return promauto.With(*r.registerer).NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: help,
	}, labels).WithLabelValues(labels...)
}

func (r *prometheusRegistry) Gauge(name, help string, labels ...string) Gauge {
	return promauto.With(*r.registerer).NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}, labels).WithLabelValues(labels...)
}

func (r *prometheusRegistry) Histogram(name, help string, labels ...string) Histogram {
	// A default set of buckets for HTTP latencies is a good starting point.
	buckets := []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5}
	return promauto.With(*r.registerer).NewHistogramVec(prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: buckets,
	}, labels).WithLabelValues(labels...)
}
