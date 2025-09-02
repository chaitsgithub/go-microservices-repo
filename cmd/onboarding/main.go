package main

import (
	"context"
	"net/http"

	health "chaits.org/go-microservices-repo/internal/handlers"
	handlers "chaits.org/go-microservices-repo/internal/handlers/onboarding"
	"chaits.org/go-microservices-repo/internal/repositories"
	appserver "chaits.org/go-microservices-repo/internal/server"
	"chaits.org/go-microservices-repo/pkg/general/logger"
	"chaits.org/go-microservices-repo/pkg/general/tracing"
	"chaits.org/go-microservices-repo/pkg/network/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const serviceName = "onboarding"

func main() {
	logger.Init(serviceName)
	shutdownTracer := tracing.InitTracer(context.Background(), serviceName)
	defer shutdownTracer()

	repos, err := repositories.NewMySQLDBManager()
	if err != nil {
		logger.Logger.WithError(err).Error("error getting DB manager")
	}

	middlewares := middleware.AllMiddlewareManager(serviceName, repos.AppRepo)
	appsHandler := handlers.NewAppsHandler(repos)

	http.Handle("/apps/list", middlewares.Then(appsHandler.GetAppsHandler, "getapps-handler"))
	http.Handle("/apps/create", middlewares.Then(appsHandler.RegisterAppHandler, "register-app-handler"))
	http.Handle("/apps/delete", middlewares.Then(appsHandler.RevokeAppHandler, "revoke-app-handler"))
	http.Handle("/health", health.HealthHandler(serviceName))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{Addr: ":8080"}
	appserver.StartServer(serviceName, server)
}
