package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shubhbham/BankingApi_Golang/internal/db"
)

type HealthHandler struct {
	db *db.DB
}

func NewHealthHandler(database *db.DB) *HealthHandler {
	return &HealthHandler{db: database}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	status := "healthy"
	statusCode := http.StatusOK

	if err := h.db.Health(r.Context()); err != nil {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"status": status,
	})
}