package app

import "net/http"

func GetRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(`/{id}/`, getLink)
	mux.HandleFunc(`/`, postLink)
	return mux
}
