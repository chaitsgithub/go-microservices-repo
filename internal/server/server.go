package appserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(serviceName string, server *http.Server) {
	go func() {
		log.Printf("%s service starting on port 8080\n", serviceName)
		server.ListenAndServe()
	}()
	gracefulShutdown(server, 10*time.Second)
}

// GracefulShutdown is a generic function to perform graceful shutdown on a running server.
// It listens for OS signals to trigger the shutdown process.
func gracefulShutdown(server *http.Server, timeout time.Duration) {
	// Create a channel to listen for OS signals.
	stop := make(chan os.Signal, 1)

	// Set up the signal handler. We want to catch SIGTERM and SIGINT.
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received.
	<-stop

	log.Println("Shutting down the server gracefully...")

	// Create a context with a timeout. This gives active requests a chance to finish.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Shut down the HTTP server.
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
		os.Exit(1)
	}

	// This is where you would add additional cleanup logic for other resources
	// like database connections, message queue consumers, etc.
	log.Println("Server has been shut down.")
}
