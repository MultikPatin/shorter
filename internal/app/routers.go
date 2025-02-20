package app

import (
	"github.com/go-chi/chi/v5"
	"main/internal/middleware"
	"net/http"
)

type Handlers interface {
	postLink(w http.ResponseWriter, r *http.Request)
	postJSONLink(w http.ResponseWriter, r *http.Request)
	getLink(w http.ResponseWriter, r *http.Request)
}

func GetRouters(h Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AccessLogger)

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.postLink)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.getLink)
		})
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", h.postJSONLink)
		})
	})
	return r
}
