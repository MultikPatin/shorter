package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"main/internal/models"
	"net/http"
)

const (
	textContentType = "text/plain; charset=utf-8"
	jsonContentType = "application/json"
)

type linksService interface {
	Add(ctx context.Context, shortenRequest models.ShortenRequest, host string) (string, error)
	AddBatch(ctx context.Context, shortenRequests []models.ShortenRequest) ([]string, error)
	Get(ctx context.Context, shortLink string) (string, error)
	Ping() error
}

func GetHandlers(s linksService) *Handlers {
	return &Handlers{
		linksService: s,
	}
}

type Handlers struct {
	linksService linksService
}

func (h *Handlers) getLink(w http.ResponseWriter, r *http.Request) {
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

func (h *Handlers) addLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	//
	var shorten models.ShortenRequest
	var response models.ShortenResponse
	//
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//if err = json.Unmarshal(buf.Bytes(), &shorten); err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	response.Result, err = h.linksService.AddBatch(ctx, shorten.URL, r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	//resp, err := json.Marshal(response)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *Handlers) addLink(w http.ResponseWriter, r *http.Request) {
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

	response.Result, err = h.linksService.Add(ctx, shortenRequest, r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *Handlers) addLinkInText(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.linksService.Add(ctx, string(body), r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", textContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func (h *Handlers) ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := h.linksService.Ping()
	if err != nil {
		http.Error(w, "Database not available", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", textContentType)
	w.WriteHeader(http.StatusOK)
}
