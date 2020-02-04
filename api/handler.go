package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func RegisterHandler(handler UserHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", handler.Get)
	r.Post("/", handler.Post)
	r.Post("/update", handler.Update)
	r.Post("/delete", handler.Delete)

	return r
}
