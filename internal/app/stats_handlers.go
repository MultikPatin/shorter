package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"main/internal/services"
	"net"
	"net/http"
	"strings"
)

// NewStatsHandlers creates a new StatsHandlers instance with the provided StatsService.
func NewStatsHandlers(s interfaces.StatsService) *StatsHandlers {
	return &StatsHandlers{
		statsService: s,
	}
}

// StatsHandlers implements the interfaces.StatsHandlers interface for handling statistics-related HTTP requests.
type StatsHandlers struct {
	statsService interfaces.StatsService
}

// GetMainStats handles GET requests for retrieving main statistics.
// It returns a JSON response with the statistics or an appropriate error.
// Possible HTTP statuses:
//   - 200 OK: statistics retrieved successfully
//   - 405 Method Not Allowed: request method is not GET
//   - 403 Forbidden: request is not from a trusted subnet
//   - 500 Internal Server Error: error marshaling response
func (h *StatsHandlers) GetMainStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var response models.MainStatsResponse

	ip, err := resolveIP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	mainStats, err := h.statsService.GetMainStats(ctx, ip.String())
	if err != nil {
		if errors.Is(err, services.ErrNotTrustedSubnet) {
			http.Error(w, "request from not trusted subnet", http.StatusForbidden)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response = models.MainStatsResponse(mainStats)

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", constants.JSONContentType)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// resolveIP extracts the client IP address from the request headers.
// It first checks the "X-Real-IP" header, then "X-Forwarded-For" if needed.
func resolveIP(r *http.Request) (net.IP, error) {
	ipStr := r.Header.Get("X-Real-IP")
	ip := net.ParseIP(ipStr)
	if ip == nil {
		ips := r.Header.Get("X-Forwarded-For")
		ipStrs := strings.Split(ips, ",")
		ipStr = ipStrs[0]
		ip = net.ParseIP(ipStr)
	}
	if ip == nil {
		return nil, fmt.Errorf("failed parse ip from http header")
	}
	return ip, nil
}
