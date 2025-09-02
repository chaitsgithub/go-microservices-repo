package handlers

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func ChainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Started Chain Call")

	// Create a new HTTP client that automatically injects trace context
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	for i := 0; i < 1; i++ {
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
	span.AddEvent("Received response from /hello", trace.WithAttributes(attribute.String("response.body", string(body))))
}

func slowProcess() {
	randomMilliseconds := rand.Intn(1000)
	time.Sleep(time.Duration(randomMilliseconds) * time.Millisecond)
	// time.Sleep(5 * time.Second)
}
