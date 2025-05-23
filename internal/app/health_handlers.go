package app

import (
	"main/internal/constants"
	"main/internal/interfaces"
	"net/http"
)

// HealthHandlers coordinates HTTP requests for health-related endpoints.
type HealthHandlers struct {
	healthService interfaces.HealthService // Service dependency for executing health checks.
}

// NewHealthHandlers constructs a new HealthHandlers instance wired up to a HealthService.
func NewHealthHandlers(s interfaces.HealthService) *HealthHandlers {
	return &HealthHandlers{
		healthService: s,
	}
}

// Ping responds to GET requests by invoking the health service and reporting status.
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

	w.Header().Set("content-type", constants.TextContentType)
	w.WriteHeader(http.StatusOK)
}
