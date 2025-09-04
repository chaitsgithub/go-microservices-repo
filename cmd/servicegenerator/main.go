package main

import (
	"log"
	"net/http"
)

func main() {
	// Register the handler for the service generator
	http.HandleFunc("/generate", GenerateServiceHandler)

	// Serve the UI for the generator
	http.Handle("/", http.FileServer(http.Dir("./")))

	log.Println("Service generator is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
