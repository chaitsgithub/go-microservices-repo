package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// Get the global Tracer instance
var tracer = otel.Tracer("tracer-service")

func helloHandler(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("handler.name", "hello-handler"))

	io.WriteString(w, "Hello, World!\n")
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

// initTracer initializes the OpenTelemetry tracer provider.
func initTracer() func() {
	ctx := context.Background()

	// Create an OTLP trace exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	// Create a resource with service name and instance ID
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("go-otel-demo-service"),
			attribute.String("service.version", "1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// Create a new trace provider with the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(time.Second*5)),
		sdktrace.WithResource(res),
	)

	// Register the global trace provider
	otel.SetTracerProvider(tp)

	log.Println("OpenTelemetry tracer initialized.")

	// Return a function to shutdown the tracer
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown trace provider: %v", err)
		}
	}
}

func main() {

	shutdownTracer := initTracer()
	defer shutdownTracer()

	// Use otelhttp.NewHandler to automatically create spans for incoming requests
	http.Handle("/hello", otelhttp.NewHandler(http.HandlerFunc(helloHandler), "hello-handler"))
	http.Handle("/chain", otelhttp.NewHandler(http.HandlerFunc(chainHandler), "chain-handler"))

	fmt.Println("Server starting on port 8080 with OpenTelemetry tracing...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
