package app

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
)

var inMemoryDB = NewInMemoryDB()
var ShortPre = ""

func postLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	u, err := uuid.NewRandom()
	if err != nil {
		http.Error(res, "Failed to generate UUID", http.StatusInternalServerError)
		return
	}

	var key string
	var response string

	if IsURL(ShortPre) {
		key = u.String()
		response = ShortPre + delimiter + key + delimiter
	} else {
		key = ShortPre + key
		response = urlPrefix + req.Host + delimiter + key + delimiter
	}

	_, err = inMemoryDB.AddLink(key, string(body))
	if err != nil {
		http.Error(res, "Failed to add link", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", contentType)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(response))
}

func getLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	id := req.PathValue("id")

	origin, err := inMemoryDB.GetByID(id)
	if err != nil {
		http.Error(res, "Origin not found", http.StatusNotFound)
		return
	}

	res.Header().Set("content-type", contentType)
	res.Header().Set("Location", origin)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
