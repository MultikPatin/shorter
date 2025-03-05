package app

import (
	"bytes"
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
	Add(origin string, host string) (string, error)
	Get(id string) (string, error)
}

func GetHandlers(s linksService) *Handlers {
	return &Handlers{
		linksService: s,
	}
}

type Handlers struct {
	linksService linksService
}

func (h *Handlers) postJSONLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var shorten models.ShortenRequest
	var response models.ShortenResponse

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &shorten); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.Result, err = h.linksService.Add(shorten.URL, r.Host)
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

func (h *Handlers) postLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	response, err := h.linksService.Add(string(body), r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", textContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func (h *Handlers) getLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	origin, err := h.linksService.Get(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Origin not found", http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", textContentType)
	w.Header().Set("Location", origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
