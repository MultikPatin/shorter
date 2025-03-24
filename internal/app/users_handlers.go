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

	defer r.Body.Close()
	var shortLinks []string
	if err := json.NewDecoder(r.Body).Decode(&shortLinks); err != nil {
		http.Error(w, "Не удалось распарсить тело запроса", http.StatusBadRequest)
		return
	}

	if len(shortLinks) == 0 {
		http.Error(w, "Список ссылок не предоставлен", http.StatusBadRequest)
		return
	}

	err := h.usersService.DeleteLinks(ctx, shortLinks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", constants.TextContentType)
	w.WriteHeader(http.StatusAccepted)
}
