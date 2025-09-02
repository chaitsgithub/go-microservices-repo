package middleware

import (
	"net/http"
	"time"

	"chaits.org/go-microservices-repo/pkg/general/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func WithPrometheusMetrics(servicename string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := newLoggingResponseWriter(w)
			next.ServeHTTP(lrw, r)
			duration := time.Since(start).Milliseconds()

			// Pass the servicename to the metrics function
			metrics.RecordHTTPRequest(servicename, r.Method, r.URL.Path, lrw.statusCode, time.Duration(duration)*time.Millisecond)
		})
	}
}

// Handler that exposes Prometheus metrics.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
