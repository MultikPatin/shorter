package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Headers interface {
	postLink(w http.ResponseWriter, r *http.Request)
	getLink(w http.ResponseWriter, r *http.Request)
}

func GetRouter(h Headers) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.postLink)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.getLink)
		})
	})
	return r
}
