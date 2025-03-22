package app

import (
	"main/internal/interfaces"
	"net/http"
)

type HealthHandlers struct {
	healthService interfaces.HealthService
}

func NewHealthHandlers(s interfaces.HealthService) *HealthHandlers {
	return &HealthHandlers{
		healthService: s,
	}
}

func (h *HealthHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	err := h.healthService.Ping()
	if err != nil {
		http.Error(w, "Database not available", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", textContentType)
	w.WriteHeader(http.StatusOK)
}
