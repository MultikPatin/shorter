package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"main/internal/interfaces"
	"main/internal/models"
	"main/internal/services"
	"net/http"
)

func NewLinksHandlers(s interfaces.LinksService) *LinksHandlers {
	return &LinksHandlers{
		linksService: s,
	}
}

type LinksHandlers struct {
	linksService interfaces.LinksService
}

func (h *LinksHandlers) GetLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	originLink, err := h.linksService.Get(ctx, r.PathValue("id"))
	if err != nil {
		http.Error(w, "Origin not found", http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", textContentType)
	w.Header().Set("Location", originLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *LinksHandlers) AddLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var shortenRequests []models.ShortensRequest
	var responses []models.ShortensResponse

	var buf bytes.Buffer
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

	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *LinksHandlers) AddLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var shortenRequest models.ShortenRequest
	var response models.ShortenResponse

	var buf bytes.Buffer
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

	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(status)
	w.Write(resp)
}

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

	w.Header().Set("content-type", textContentType)
	w.WriteHeader(status)
	w.Write([]byte(response))
}
