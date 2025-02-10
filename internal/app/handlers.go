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

	key := getDBKey(u, ShortPre)
	response := getResponseLink(key, ShortPre, urlPrefix+r.Host)

	_, err = inMemoryDB.AddLink(key, string(body))
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
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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

func getDBKey(u uuid.UUID, p string) string {
	if IsURL(p) {
		return u.String()
	}
	return p + u.String()
}

func getResponseLink(k string, p string, h string) string {
	if IsURL(p) {
		return p + delimiter + k + delimiter
	}
	return h + delimiter + k + delimiter
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
