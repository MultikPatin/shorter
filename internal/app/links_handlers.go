package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"main/internal/services"
	"net/http"
)

// NewLinksHandlers constructs a new LinksHandlers instance initialized with a LinksService.
func NewLinksHandlers(s interfaces.LinksService) *LinksHandlers {
	return &LinksHandlers{
		linksService: s,
	}
}

// LinksHandlers encapsulates handlers for managing links and redirections.
type LinksHandlers struct {
	linksService interfaces.LinksService // Dependency injection of the links service.
}

// GetLink handles GET requests for resolving short links to their original URLs.
//
// Possible HTTP statuses:
//   - 200 OK: Successfully redirected to the original URL.
//   - 404 Not Found: Original URL was not found.
//   - 410 Gone: Original URL has been deleted.
//   - 405 Method Not Allowed: Request method is not allowed (only GET supported).
func (h *LinksHandlers) GetLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	originLink, err := h.linksService.Get(ctx, r.PathValue("id"))
	if err != nil {
		if errors.Is(err, services.ErrDeletedLink) {
			http.Error(w, "Origin is deleted", http.StatusGone)
		} else {
			http.Error(w, "Origin not found", http.StatusNotFound)
		}
		return
	}

	w.Header().Set("content-type", constants.TextContentType)
	w.Header().Set("Location", originLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// AddLinks processes POST requests for batch-link creation.
//
// Possible HTTP statuses:
//   - 201 Created: All links successfully created.
//   - 400 Bad Request: Malformed request body.
//   - 405 Method Not Allowed: Request method is not allowed (only POST supported).
//   - 500 Internal Server Error: An internal error occurred during link creation.
func (h *LinksHandlers) AddLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var shortenRequests []models.ShortensRequest
	var responses []models.ShortensResponse

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &shortenRequests); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var originLinks []models.OriginLink
	for _, req := range shortenRequests {
		originLink := models.OriginLink(req)
		originLinks = append(originLinks, originLink)
	}

	results, err := h.linksService.AddBatch(ctx, originLinks, r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, result := range results {
		shortensResponse := models.ShortensResponse(result)
		responses = append(responses, shortensResponse)
	}

	resp, err := json.Marshal(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", constants.JSONContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// AddLink handles individual link creation via POST requests.
//
// Possible HTTP statuses:
//   - 201 Created: Link successfully created.
//   - 400 Bad Request: Malformed request body.
//   - 405 Method Not Allowed: Request method is not allowed (only POST supported).
//   - 409 Conflict: Duplicate link already exists.
//   - 500 Internal Server Error: An internal error occurred during link creation.
func (h *LinksHandlers) AddLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var shortenRequest models.ShortenRequest
	var response models.ShortenResponse

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &shortenRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	originLink := models.OriginLink{
		URL: shortenRequest.URL,
	}

	status := http.StatusCreated

	response.Result, err = h.linksService.Add(ctx, originLink, r.Host)
	if err != nil {
		if errors.Is(err, services.ErrConflict) {
			status = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", constants.JSONContentType)
	w.WriteHeader(status)
	w.Write(resp)
}

// AddLinkInText processes link creation directly from plain-text bodies.
//
// Possible HTTP statuses:
//   - 201 Created: Link successfully created.
//   - 405 Method Not Allowed: Request method is not allowed (only POST supported).
//   - 409 Conflict: Duplicate link already exists.
//   - 500 Internal Server Error: An internal error occurred during link creation.
func (h *LinksHandlers) AddLinkInText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	originLink := models.OriginLink{
		URL: string(body),
	}

	status := http.StatusCreated

	response, err := h.linksService.Add(ctx, originLink, r.Host)
	if err != nil {
		if errors.Is(err, services.ErrConflict) {
			status = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("content-type", constants.TextContentType)
	w.WriteHeader(status)
	w.Write([]byte(response))
}
