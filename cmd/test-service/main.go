package main

import (
	"context"
	"net/http"
	"time"

	health "chaits.org/go-microservices-repo/internal/handlers"
	handlers "chaits.org/go-microservices-repo/internal/handlers/test-service"
	"chaits.org/go-microservices-repo/internal/repositories"
	appserver "chaits.org/go-microservices-repo/internal/server"
	"chaits.org/go-microservices-repo/pkg/general/logger"
	"chaits.org/go-microservices-repo/pkg/general/tracing"
	"chaits.org/go-microservices-repo/pkg/network/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var serviceName = "test-service"

func main() {

	logger.Init(serviceName)

	shutdownTracer := tracing.InitTracer(context.Background(), serviceName)
	defer shutdownTracer()

	repos, err := repositories.NewMySQLDBManager()
	if err != nil {
		logger.Logger.WithError(err).Fatal("DB Error")
	}

	middlewares := middleware.NewManager(
		middleware.WithLogging,
		middleware.WithPrometheusMetrics(serviceName),
		middleware.WithCORS,
		middleware.WithRateLimiter(100, time.Minute),
		middleware.WithAPIKeyAuth(repos.AppRepo),
	)

	http.Handle("/hello", middlewares.Then(handlers.HelloHandler, "hello-handler"))
	http.Handle("/chain", middlewares.Then(handlers.ChainHandler, "chain-handler"))
	http.Handle("/health", health.HealthHandler(serviceName))

	// http.Handle("/hello", otelhttp.NewHandler(middlewares.Then1(handlers.HelloHandler), "hello-handler"))
	// http.Handle("/chain", otelhttp.NewHandler(middlewares.Then1(handlers.ChainHandler), "chain-handler"))

	// http.Handle("/hello", otelhttp.NewHandler(middleware.ChainAllHandlers(handlers.HelloHandler, serviceName), "hello-handler"))
	// http.Handle("/chain", otelhttp.NewHandler(middleware.ChainAllHandlers(handlers.ChainHandler, serviceName), "chain-handler"))

	server := &http.Server{Addr: ":8081"}
	http.Handle("/metrics", promhttp.Handler())
	appserver.StartServer(serviceName, server)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

func chiRouter() {
	// r := chi.NewRouter()
	// r.Use(middleware.WithPrometheusMetrics, middleware.WithLogging)
	// r.Method(http.MethodGet, "/hello", otelhttp.NewHandler(middleware.ChainAllHandlers(helloHandler), "hello-handler"))
	// r.Method(http.MethodGet, "/chain", otelhttp.NewHandler(middleware.ChainAllHandlers(chainHandler), "chain-handler"))
	// http.Handle("/", r)
	// http.Handle("/metrics", promhttp.Handler())
	// 	fmt.Println("Server starting on port 8080 with OpenTelemetry tracing...")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
