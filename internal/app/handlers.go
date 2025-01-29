package app

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
)

const (
	urlPrefix   = "http://"
	delimiter   = "/"
	contentType = "text/plain; charset=utf-8"
)

var db = NewInMemoryDB()

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
	id, err := db.AddLink(u, string(body))
	if err != nil {
		http.Error(res, "Failed to add link", http.StatusInternalServerError)
		return
	}

	response := urlPrefix + req.Host + delimiter + id.String() + delimiter

	res.Header().Set("content-type", contentType)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(response))
}

func getLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	short := req.PathValue("id")
	id, err := uuid.Parse(short)
	if err != nil {
		http.Error(res, "Path value is not valid UUID", http.StatusBadRequest)
		return
	}
	origin, err := db.GetByID(id)
	if err != nil {
		http.Error(res, "Origin not found", http.StatusNotFound)
		return
	}
	fmt.Println(origin)

	res.Header().Set("content-type", contentType)
	res.Header().Set("Location", origin)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
