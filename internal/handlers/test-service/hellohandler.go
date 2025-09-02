package handlers

import (
	"context"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("handler.name", "hello-handler"))

	// Start a new span to simulate a DB call
	dbCtx, dbSpan := otel.Tracer("db-tracer").Start(ctx, "db.query")
	defer dbSpan.End()
	// Simulate DB operation
	performDBQuery(dbCtx)
	time.Sleep(20 * time.Millisecond)
	dbSpan.SetStatus(codes.Ok, "DB query successful")
	dbSpan.End()

	// Start another span for a cache lookup
	cacheCtx, cacheSpan := otel.Tracer("cache-tracer").Start(ctx, "cache.lookup")
	defer cacheSpan.End()
	// Simulate cache lookup
	time.Sleep(5 * time.Millisecond)
	performCacheLookup(cacheCtx)
	cacheSpan.SetStatus(codes.Ok, "Cache hit")
	cacheSpan.End()

	slowProcess()
	io.WriteString(w, "Hello, World!\n")
}

// performDBQuery simulates a database query.
func performDBQuery(ctx context.Context) {
	// Now, any nested spans or operations here will be children of the "db.query" span.
	// For example, you could add another span for a specific table query:
	_, tableSpan := otel.Tracer("db-tracer").Start(ctx, "db.query.users_table")
	defer tableSpan.End()

	// Simulate the actual database operation
	time.Sleep(20 * time.Millisecond)
}

// performCacheLookup simulates a cache lookup operation.
func performCacheLookup(ctx context.Context) {
	// We are using the context with the "cache.lookup" span.
	// This allows us to add events or attributes specific to the cache operation.
	span := trace.SpanFromContext(ctx)
	span.AddEvent("cache.check")

	_, tableSpan := otel.Tracer("cache-tracer").Start(ctx, "cache.check")
	defer tableSpan.End()

	// Simulate the cache lookup time.
	time.Sleep(5 * time.Millisecond)
}
