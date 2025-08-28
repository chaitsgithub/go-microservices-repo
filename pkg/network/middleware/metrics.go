package middleware

import (
	"net/http"
	"time"

	"chaits.org/go-microservices-repo/pkg/general/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func WithPrometheusMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Wrap the response writer to capture the status code.
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		duration := time.Since(start).Seconds()
		// Record Metrics
		metrics.RecordHTTPRequest(r.Method, r.URL.Path, lrw.statusCode, time.Duration(duration))
	})
}

// Handler that exposes Prometheus metrics.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
