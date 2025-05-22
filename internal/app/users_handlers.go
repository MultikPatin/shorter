package app // Package app implements HTTP handlers for user-related operations.

import (
	"encoding/json"
	"errors"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"main/internal/services"
	"net/http"
)

// NewUsersHandlers constructs a new UsersHandlers instance injected with a UsersService.
func NewUsersHandlers(s interfaces.UsersService) *UsersHandlers {
	return &UsersHandlers{
		usersService: s,
	}
}

// UsersHandlers encapsulates handlers for user-specific operations like fetching and deleting links.
type UsersHandlers struct {
	usersService interfaces.UsersService // Dependency injection of the users service.
}

// GetLinks handles GET requests for retrieving links associated with the currently-authenticated user.
//
// Possible HTTP statuses:
//   - 200 OK: Successfully fetched the user's links.
//   - 204 No Content: The user has no links.
//   - 405 Method Not Allowed: Request method is not allowed (only GET supported).
//   - 500 Internal Server Error: An internal error occurred during link retrieval.
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

// DeleteLinks processes DELETE requests for removing user-owned links.
//
// Possible HTTP statuses:
//   - 202 Accepted: Deletion initiated successfully.
//   - 400 Bad Request: Invalid or missing request body.
//   - 500 Internal Server Error: An internal error occurred during link deletion.
func (h *UsersHandlers) DeleteLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	var shortLinks []string
	if err := json.NewDecoder(r.Body).Decode(&shortLinks); err != nil {
		http.Error(w, "Couldn't parse the request body", http.StatusBadRequest)
		return
	}

	if len(shortLinks) == 0 {
		http.Error(w, "The list of links is not provided", http.StatusBadRequest)
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
