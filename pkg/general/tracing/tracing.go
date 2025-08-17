package tracing

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func InitTracer(ctx context.Context, serviceName string) func() {

	// Create an OTLP trace exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	// Create a resource with service name and instance ID
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
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
