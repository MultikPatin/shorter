package app

import (
	"github.com/go-chi/chi/v5"
	"main/internal/middleware"
)

func NewRouters(h *Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AccessLogger)
	r.Use(middleware.GZipper)
	r.Use(middleware.Authentication)

	r.Route("/", func(r chi.Router) {
		r.Post("/login", h.users.Login)
		r.Get("/ping", h.health.Ping)
		r.Post("/", h.links.AddLinkInText)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.links.GetLink)
		})
		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", h.links.AddLink)
				r.Post("/batch", h.links.AddLinks)
			})
		})
	})
	return r
}
