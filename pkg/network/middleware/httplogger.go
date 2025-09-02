// Filename: network/middleware/logging_middleware.go
// Package: middleware
// Description: HTTP middleware to log detailed request and response information.

package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"chaits.org/go-microservices-repo/pkg/general/logger"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// loggingResponseWriter is a custom wrapper around http.ResponseWriter to
// capture the response status code and body.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

// newLoggingResponseWriter creates a new instance of our custom writer.
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
		body:           new(bytes.Buffer),
	}
}

// WriteHeader captures the status code.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Write captures the response body and passes it through to the original writer.
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

// getOtelTraceID retrieves the OpenTelemetry trace ID from the request context.
// This is the correct way to get the ID that's propagated across services.
func getOtelTraceID(ctx context.Context) string {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return "N/A"
}

// WithLogging is the middleware function. It wraps an http.Handler and
// logs detailed request and response information.
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Read the request body. We create a copy so the handler can still read it.
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.WithError(err).Error("Failed to read request body for logging.")
		}
		// Put the body back into the request.
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		// Create our custom logging writer.
		lrw := newLoggingResponseWriter(w)

		// Call the next handler in the chain.
		next.ServeHTTP(lrw, r)
		// Get the trace ID from the request context.
		traceID := getOtelTraceID(r.Context())

		// Log all the captured details after the request is finished.
		duration := time.Since(start)
		logger.Logger.WithFields(logrus.Fields{
			"traceID":       traceID,
			"method":        r.Method,
			"path":          r.URL.Path,
			"query":         r.URL.RawQuery,
			"request_body":  string(reqBody),
			"status_code":   lrw.statusCode,
			"response_body": lrw.body.String(),
			"duration_ms":   duration.Milliseconds(),
		}).Info("HTTP request processed.")
	})
}
