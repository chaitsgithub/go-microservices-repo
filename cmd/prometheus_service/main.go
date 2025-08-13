package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"chaits.org/microservices-repo/pkg/general/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ServiceMetrics holds all the metrics for this service.
// This is what our microservice will depend on.
type ServiceMetrics struct {
	requestCounter    metrics.Counter
	activeConnections metrics.Gauge
	requestLatency    metrics.Histogram
}

// NewServiceMetrics initializes and returns all the metrics for the service.
func NewServiceMetrics(registry metrics.Registry) *ServiceMetrics {
	return &ServiceMetrics{
		requestCounter:    registry.Counter("http_requests_total", "Total number of HTTP requests."),
		activeConnections: registry.Gauge("http_active_connections", "Number of currently active HTTP connections."),
		requestLatency:    registry.Histogram("http_request_duration_seconds", "HTTP request latency in seconds."),
	}
}

// HelloHandler is a sample HTTP handler that uses our metrics.
func HelloHandler(m *ServiceMetrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Increment the counter for every request.
		m.requestCounter.Inc()

		// Simulate some work with a random delay to test the histogram.
		start := time.Now()
		time.Sleep(time.Duration(100+r.ContentLength%500) * time.Millisecond) // Simulate work
		m.requestLatency.Observe(time.Since(start).Seconds())

		// Simulate an active connection using the gauge.
		m.activeConnections.Inc()
		defer m.activeConnections.Dec()

		fmt.Fprintln(w, "Hello, world!")
	}
}

func main() {
	metricsBackend := metrics.METRICS_PROMETHEUS
	prometheusRegistry, err := metrics.NewMetricsRegistry(metricsBackend)
	if err != nil {
		log.Printf("Error connecting to metrics registry. Error : %v/n", err)
	}
	serviceMetrics := NewServiceMetrics(prometheusRegistry)

	// 4. Set up HTTP server routes.
	http.Handle("/hello", HelloHandler(serviceMetrics))
	http.Handle("/metrics", promhttp.Handler()) // Expose Prometheus metrics

	log.Printf("Starting server on :8080 with '%s' metrics backend", metricsBackend)
	log.Println("Metrics available at http://localhost:8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
