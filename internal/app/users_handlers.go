package app

import (
	"encoding/json"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"net/http"
)

func NewUsersHandlers(s interfaces.UsersService) *UsersHandlers {
	return &UsersHandlers{
		usersService: s,
	}
}

type UsersHandlers struct {
	usersService interfaces.UsersService
}

func (h *UsersHandlers) GetLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var responses []models.UserLinksResponse

	results, err := h.usersService.GetLinks(ctx)
	if err != nil {
		http.Error(w, "Links not found", http.StatusInternalServerError)
		return
	}

	for _, result := range results {
		shortensResponse := models.UserLinksResponse(result)
		responses = append(responses, shortensResponse)
	}

	resp, err := json.Marshal(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", constants.JSONContentType)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
