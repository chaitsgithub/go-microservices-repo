package middleware

import (
	"net/http"
)

// Chain Usage:
// Chain (handler, middleware1, middleware2, middleware3)
// The middleware will execute in the order: middleware1 -> middleware2 -> middleware3 -> handler
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// Loop through the middleware slice in reverse order to build the chain.
	for i := len(middlewares) - 1; i >= 0; i-- {
		// Wrap the current handler with the next middleware in the chain.
		h = middlewares[i](h)
	}
	return h
}

func ChainAllHandlers(h http.HandlerFunc) http.Handler {
	return Chain(http.HandlerFunc(h), WithLogging, WithPrometheusMetrics)
}
