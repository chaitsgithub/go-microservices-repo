// Package metrics to record Prometheus Metrics
package metrics

import (
	"strconv"
	"time"
)

// --- Request-Level Metrics Utilities ---

// RecordHTTPRequest records a single HTTP request, updating the total count,
// duration, and error count (if applicable).
// It should be called after a request has been handled.
func RecordHTTPRequest(serviceName, method, path string, statusCode int, duration time.Duration) {
	status := strconv.Itoa(statusCode)
	HTTPRequestTotal.WithLabelValues(serviceName, method, path, status).Inc()
	HTTPRequestDurationSeconds.WithLabelValues(serviceName, method, path, status).Observe(duration.Seconds())

	// Increment the error counter if the status code indicates an error (5xx).
	if statusCode >= 500 && statusCode < 600 {
		HTTPRequestsErrorsTotal.WithLabelValues(serviceName, method, path).Inc()
	}
}

// --- Dependency-Level Metrics Utilities ---

// RecordDependencyRequest measures the duration and records a single request to an external dependency.
// It should be called after a dependency call is completed.
func RecordDependencyRequest(dependencyName, status string, duration time.Duration) {
	DependencyRequestsTotal.WithLabelValues(dependencyName, status).Inc()
	DependencyDurationSeconds.WithLabelValues(dependencyName, status).Observe(duration.Seconds())
}

// UpdateDatabaseConnections sets the value of the database connections gauge.
// Call this function periodically to report the number of open connections.
func UpdateDatabaseConnections(count int) {
	DatabaseConnectionsOpen.Set(float64(count))
}

// --- Application-Specific Metrics Utilities ---

// UpdateJobQueueSize sets the value of the job queue size gauge.
// This function can be called periodically to report the current queue length.
func UpdateJobQueueSize(size int) {
	JobQueueSize.Set(float64(size))
}

// IncrementUserRegistrations increments the total count of new user registrations.
// Call this function whenever a new user signs up.
func IncrementUserRegistrations() {
	UserRegistrationsTotal.Inc()
}

// IncrementCheckoutEvents increments the total count of completed checkout events.
// Call this function at the successful completion of a checkout.
func IncrementCheckoutEvents() {
	CheckoutEventsTotal.Inc()
}
