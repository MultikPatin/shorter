package app

import (
	"github.com/go-chi/chi/v5"
	"main/internal/middleware"
)

// NewRouters constructs and configures the main router with middleware and routes.
func NewRouters(h *Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AccessLogger)
	r.Use(middleware.GZipper)
	r.Use(middleware.Authentication)

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", h.health.Ping)
		r.Post("/", h.links.AddLinkInText)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.links.GetLink)
		})
		r.Route("/api", func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {
				r.Get("/urls", h.users.GetLinks)
				r.Delete("/urls", h.users.DeleteLinks)
			})
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", h.links.AddLink)
				r.Post("/batch", h.links.AddLinks)
			})
		})
	})
	return r
}
