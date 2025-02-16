package app

import (
	"github.com/google/uuid"
	"io"
	"main/internal/services"
	"net/http"
)

const (
	urlPrefix   = "http://"
	contentType = "text/plain; charset=utf-8"
)

type DB interface {
	AddLink(id string, link string) (string, error)
	GetByID(id string) (string, error)
}

func GetHandlers(db DB) *MyHandlers {
	return &MyHandlers{
		database: db,
	}
}

type MyHandlers struct {
	database DB
}

var ShortPre = ""

func (h *MyHandlers) postLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	u, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, "Failed to generate UUID", http.StatusInternalServerError)
		return
	}

	key := services.GetDBKey(u, ShortPre)
	response := services.GetResponseLink(key, ShortPre, urlPrefix+r.Host)

	_, err = h.database.AddLink(key, string(body))
	if err != nil {
		http.Error(w, "Failed to add link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", contentType)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func (h *MyHandlers) getLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	origin, err := h.database.GetByID(id)
	if err != nil {
		http.Error(w, "Origin not found", http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", contentType)
	w.Header().Set("Location", origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
