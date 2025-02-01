package app

import (
	"github.com/go-chi/chi/v5"
)

func GetRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", postLink)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getLink)
		})
	})
	return r
}
