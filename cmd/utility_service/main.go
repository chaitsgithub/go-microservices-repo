// cmd/hello_service/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Define an HTTP handler
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, "Hello World\n")
		} else {
			http.Error(w, "Method not Supported", http.StatusMethodNotAllowed)
		}
	})

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listenAddr := fmt.Sprintf(":%s", port)
	log.Printf("Starting hello_service on %s\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
