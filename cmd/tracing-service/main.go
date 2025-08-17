package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"chaits.org/microservices-repo/pkg/general/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var serviceName = "tracing-service"

func helloHandler(w http.ResponseWriter, r *http.Request) {
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

func slowProcess() {
	time.Sleep(100 * time.Millisecond)
}

func chainHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Started Chain Call")

	// Create a new HTTP client that automatically injects trace context
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	for i := 0; i < 3; i++ {
		callHello(ctx, span, client, w)
		// Simulate a slow process
		slowProcess()
		io.WriteString(w, "Chained call complete!\n")
	}
}

func callHello(ctx context.Context, span trace.Span, client http.Client, w http.ResponseWriter) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/hello", nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create Request")
		http.Error(w, fmt.Sprintf("Error Creating Request with Context: %w", err), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed call to /hello")
		http.Error(w, fmt.Sprintf("Error calling /hello: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Received response from /hello: %s", string(body))
	span.AddEvent("Received response from /hello", trace.WithAttributes(attribute.String("response.body", string(body))))
}

func main() {

	shutdownTracer := tracing.InitTracer(context.Background(), serviceName)
	defer shutdownTracer()

	// Use otelhttp.NewHandler to automatically create spans for incoming requests
	http.Handle("/hello", otelhttp.NewHandler(http.HandlerFunc(helloHandler), "hello-handler"))
	http.Handle("/chain", otelhttp.NewHandler(http.HandlerFunc(chainHandler), "chain-handler"))

	fmt.Println("Server starting on port 8080 with OpenTelemetry tracing...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
