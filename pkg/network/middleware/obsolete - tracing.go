package middleware

// Obsolete Code when using OTEL

// import (
// 	"context"
// 	"log"
// 	"net/http"

// 	"chaits.org/go-microservices-repo/pkg/general/logger"
// 	"github.com/google/uuid"
// )

// type traceIDKey string

// const (
// 	// TraceIDHeader is the name of the HTTP header used to pass the trace ID.
// 	TraceIDHeader = "X-Trace-ID"
// 	// traceIDContextKey is the key used to store the trace ID in the request context.
// 	traceIDContextKey traceIDKey = "traceID"
// )

// // WithTracing is an HTTP middleware that adds a trace ID to the request context.
// // It checks for an existing trace ID in the "X-Trace-ID" header. If none is found,
// // it generates a new UUID. This is essential for tracing requests across microservices.
// func WithTracing(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Get or generate a trace ID for the request.
// 		logger.Logger.Infof("Request Context : ", r.Context())
// 		traceID, _ := GetTraceIDFromContext(r.Context())
// 		if traceID == "" {
// 			traceID = r.Header.Get(TraceIDHeader)
// 			if traceID == "" {
// 				traceID = uuid.New().String()
// 			}
// 		}

// 		// Add the trace ID to the request context.
// 		ctx := context.WithValue(r.Context(), traceIDContextKey, traceID)
// 		r = r.WithContext(ctx)

// 		// Log the incoming request with the trace ID.
// 		log.Printf("Received request. TraceID: %s, Method: %s, URI: %s", traceID, r.Method, r.RequestURI)

// 		// Call the next handler in the chain.
// 		next.ServeHTTP(w, r)
// 	})
// }

// // GetTraceIDFromContext is a helper function to retrieve the trace ID from the request context.
// // It returns the trace ID as a string and a boolean indicating if the key was found.
// func GetTraceIDFromContext(ctx context.Context) (string, bool) {
// 	traceID, ok := ctx.Value(traceIDContextKey).(string)
// 	return traceID, ok
// }
