// Filename: metrics/metrics.go
// Package: metrics
// Description: Centralized definitions for common Prometheus metrics for microservices.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// --- 1. Request-Level Metrics (RED Method) ---
var (
	// HttpRequestTotal is a CounterVec to count total HTTP requests.
	// Labels differentiate requests by method, path, and status code.
	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"serviceName", "method", "path", "status_code"},
	)

	// HttpRequestDurationSeconds is a HistogramVec to measure request duration.
	// This uses a default set of buckets for common web latency ranges.
	HTTPRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_requests_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"serviceName", "method", "path", "status_code"},
	)

	// HttpRequestsErrorsTotal is a CounterVec for HTTP requests that result in errors.
	// This helps track the number of failed requests.
	HTTPRequestsErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_errors_total",
			Help: "Total number of HTTP requests that resulted in an error.",
		},
		[]string{"serviceName", "method", "path"},
	)
)

// --- 2. Resource-Level Metrics ---
var (
	// GoGoroutinesTotal is a Gauge that tracks the current number of goroutines.
	// This metric is automatically collected by the Go client library.
	GoGoroutinesTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_goroutines_total",
			Help: "Total number of goroutines that currently exist.",
		},
	)

	// GoMemoryUsageBytes is a Gauge that tracks the current memory usage of the process.
	// This metric is automatically collected by the Go client library.
	GoMemoryUsageBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_memory_usage_bytes",
			Help: "Current memory usage of the process.",
		},
	)

	// CpuUsageSecondsTotal is a Counter that tracks the total CPU time consumed.
	// This metric is automatically collected by the Go client library.
	CPUUsageSecondsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cpu_usage_seconds_total",
			Help: "Total CPU time consumed by the process.",
		},
	)
)

// --- 3. Dependency-Level Metrics ---
var (
	// DependencyRequestsTotal is a CounterVec for requests made to external dependencies.
	// This helps monitor the health and traffic to external services.
	DependencyRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dependency_requests_total",
			Help: "Total number of requests to external dependencies.",
		},
		[]string{"dependency_name", "status"},
	)

	// DependencyDurationSeconds is a HistogramVec to measure the duration of dependency calls.
	// Essential for identifying slow dependencies.
	DependencyDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dependency_duration_seconds",
			Help:    "Duration of requests to external dependencies in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"dependency_name", "status"},
	)

	// DatabaseConnectionsOpen is a Gauge for the number of open database connections.
	// This helps manage connection pools.
	DatabaseConnectionsOpen = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_open",
			Help: "Number of open database connections.",
		},
	)
)

// --- 4. Application-Specific Metrics ---
var (
	// UserRegistrationsTotal is a Counter for the total number of new user registrations.
	// An example of a key business metric.
	UserRegistrationsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of new user registrations.",
		},
	)

	// CheckoutEventsTotal is a Counter for completed checkout events.
	// Another example of a business-specific metric.
	CheckoutEventsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "checkout_events_total",
			Help: "Total number of completed checkout events.",
		},
	)

	// JobQueueSize is a Gauge for the number of pending jobs in a queue.
	// Useful for monitoring the backlog of work.
	JobQueueSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "job_queue_size",
			Help: "Number of pending items in the processing queue.",
		},
	)
)

// init registers all defined metrics with the default Prometheus registry.
// This function is automatically called when the package is imported.
func init() {
	prometheus.MustRegister(
		HTTPRequestTotal,
		HTTPRequestDurationSeconds,
		HTTPRequestsErrorsTotal,
		GoGoroutinesTotal,
		GoMemoryUsageBytes,
		CPUUsageSecondsTotal,
		DependencyRequestsTotal,
		DependencyDurationSeconds,
		DatabaseConnectionsOpen,
		UserRegistrationsTotal,
		CheckoutEventsTotal,
		JobQueueSize,
	)
}
