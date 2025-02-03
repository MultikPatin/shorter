package app

import (
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
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
	id, err := inMemoryDB.AddLink(delimiter+u.String()+delimiter, string(body))
	if err != nil {
		http.Error(res, "Failed to add link", http.StatusInternalServerError)
		return
	}

	var response string

	log.Printf(u.String())
	log.Printf(EnvConfig.ShorLink)
	log.Printf(CmdConfig.ShorLink.Addr)

	switch {
	case EnvConfig.ShorLink != "":
		response = urlPrefix + req.Host + delimiter + EnvConfig.ShorLink + id
	case CmdConfig.ShorLink.Addr != "":
		response = urlPrefix + req.Host + delimiter + CmdConfig.ShorLink.Addr + id
	default:
		response = urlPrefix + req.Host + delimiter + id + delimiter
	}

	log.Printf(response)

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

	switch {
	case EnvConfig.ShorLink != "":
		id = strings.TrimPrefix(id, delimiter+EnvConfig.ShorLink)
	case CmdConfig.ShorLink.Addr != "":
		id = strings.TrimPrefix(id, delimiter+CmdConfig.ShorLink.Addr)
	}

	log.Printf("ID %s", id)

	origin, err := inMemoryDB.GetByID(id)
	if err != nil {
		http.Error(res, "Origin not found", http.StatusNotFound)
		return
	}

	res.Header().Set("content-type", contentType)
	res.Header().Set("Location", origin)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

//func IsUrl(str string) bool {
//	u, err := url.Parse(str)
//	return err == nil && u.Scheme != "" && u.Host != ""
//}
