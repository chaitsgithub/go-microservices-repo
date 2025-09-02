package middleware

import (
	"log"
	"net/http"

	"chaits.org/go-microservices-repo/internal/repositories"
)

// WithAPIKeyAuth is a middleware that validates an API key from a request header
// by checking it against the database.
func WithAPIKeyAuth(appRepo repositories.AppRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the API key from the X-API-Key header.
			appName := r.Header.Get("X-App-Name")
			apiKey := r.Header.Get("X-API-Key")

			if apiKey == "" {
				log.Println("Unauthorized: API key is missing.")
				http.Error(w, "Unauthorized: API Key Missing", http.StatusUnauthorized)
				return
			}

			// Validate the API key using a database lookup.
			_, ok, err := appRepo.ValidateAPIKey(r.Context(), appName, apiKey)
			if err != nil {
				log.Printf("Error validating API key: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !ok {
				log.Printf("Unauthorized: Invalid API key '%s' for App '%s'", apiKey, appName)
				http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
				return
			}

			log.Printf("Authenticated: Request with valid API key")

			// If the key is valid, proceed to the next handler.
			next.ServeHTTP(w, r)
		})
	}
}
