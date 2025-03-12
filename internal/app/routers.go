package app

import (
	"github.com/go-chi/chi/v5"
	"main/internal/interfaces"
	"main/internal/middleware"
)

func NewRouters(h interfaces.LinkHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AccessLogger)
	r.Use(middleware.GZipper)

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", h.Ping)
		r.Post("/", h.AddLinkInText)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetLink)
		})
		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", h.AddLink)
				r.Post("/batch", h.AddLinks)
			})
		})
	})
	return r
}
