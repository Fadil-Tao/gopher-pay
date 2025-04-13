package rest

import (
	"encoding/json"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := map[string]string{
		"environment": "dev",
		"status":      "healthy",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}