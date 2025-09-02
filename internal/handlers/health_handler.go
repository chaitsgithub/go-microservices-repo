package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

func HealthHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := &HealthResponse{
			Service: serviceName,
			Status:  "Healthy",
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Error formatting response", http.StatusInternalServerError)
			return
		}
	}
}
