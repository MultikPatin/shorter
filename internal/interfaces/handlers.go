package interfaces

import "net/http"

type HealthHandlers interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

type LinkHandlers interface {
	AddLinkInText(w http.ResponseWriter, r *http.Request)
	AddLink(w http.ResponseWriter, r *http.Request)
	AddLinks(w http.ResponseWriter, r *http.Request)
	GetLink(w http.ResponseWriter, r *http.Request)
}

type UsersHandlers interface {
	GetLinks(w http.ResponseWriter, r *http.Request)
}
