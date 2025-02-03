package app

import (
	"github.com/google/uuid"
	"io"
	"net/http"
)

var inMemoryDB = NewInMemoryDB()

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

	keyPrefix := EnvConfig.ShorLink
	if keyPrefix == "" {
		keyPrefix = CmdConfig.ShorLink.Addr
	}

	id, err := inMemoryDB.AddLink(keyPrefix+u.String(), string(body))
	if err != nil {
		http.Error(res, "Failed to add link", http.StatusInternalServerError)
		return
	}

	response := urlPrefix + req.Host + delimiter + id + delimiter

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
