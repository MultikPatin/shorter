package app

import (
	"encoding/json"
	"errors"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"main/internal/services"
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

	status := http.StatusOK

	results, err := h.usersService.GetLinks(ctx, r.Host)
	if err != nil {
		if errors.Is(err, services.ErrNoLinksByUser) {
			status = http.StatusNoContent
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for _, result := range results {
		response := models.UserLinksResponse(result)
		responses = append(responses, response)
	}

	resp, err := json.Marshal(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", constants.JSONContentType)
	w.WriteHeader(status)
	w.Write(resp)
}

func (h *UsersHandlers) DeleteLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var shortLinks []string

	h.usersService.DeleteLinks(ctx, shortLinks)

	w.Header().Set("content-type", constants.TextContentType)
	w.WriteHeader(http.StatusAccepted)
}
