package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Manager is a struct to manage and chain different middleware stacks.
type Manager struct {
	// A slice to store the middleware functions.
	chain []func(http.Handler) http.Handler
}

// NewManager returns a new Manager with the given middleware.
func NewManager(m ...func(http.Handler) http.Handler) *Manager {
	return &Manager{
		chain: m,
	}
}

// NewManager returns a new Manager with the given middleware.
func AllMiddlewareManager(serviceName string) *Manager {
	return NewManager(
		WithLogging,
		WithPrometheusMetrics(serviceName),
		WithCORS,
		WithRateLimiter(100, time.Minute),
	)
}

// Then adds more middleware functions to the chain.
func (m *Manager) Then(h http.HandlerFunc, otelOperation string) http.Handler {
	handler := http.Handler(h)
	chainWithOtel := append([]func(http.Handler) http.Handler{otelhttp.NewMiddleware(otelOperation)}, m.chain...)
	// Loop through the middleware slice in reverse to build the chain.
	for i := len(chainWithOtel) - 1; i >= 0; i-- {
		handler = chainWithOtel[i](handler)
	}
	return handler
}

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

func ChainAllHandlers(h http.HandlerFunc, serviceName string) http.Handler {
	return Chain(http.HandlerFunc(h), WithLogging, WithPrometheusMetrics(serviceName), WithCORS)
}
