package app

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
)

var inMemoryDB = NewInMemoryDB()
var ShortPre = ""

func postLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
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

	response := getResponseLink(u.String(), ShortPre, urlPrefix+r.Host)

	_, err = inMemoryDB.AddLink(u.String(), string(body))
	if err != nil {
		http.Error(w, "Failed to add link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", contentType)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func getLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	origin, err := inMemoryDB.GetByID(id)
	if err != nil {
		http.Error(w, "Origin not found", http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", contentType)
	w.Header().Set("Location", origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func getResponseLink(k string, p string, h string) string {
	key := delimiter + k + delimiter
	if IsURL(p) {
		return p + key
	}
	return h + p + key
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
