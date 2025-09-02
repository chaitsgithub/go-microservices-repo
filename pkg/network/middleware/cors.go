package middleware

import (
	"log"
	"net/http"
)

// WithCORS is an HTTP middleware that adds Cross-Origin Resource Sharing (CORS)
// headers to the response. It also handles preflight OPTIONS requests,
// ensuring your API can be consumed by web clients from different domains.
func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers to allow all origins, common methods, and common headers.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Log the request and the CORS headers being set.
		log.Printf("CORS middleware: Setting headers for Method: %s, Origin: %s", r.Method, r.Header.Get("Origin"))

		// Handle preflight OPTIONS requests.
		// Browsers send these to check permissions before making the actual request.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
