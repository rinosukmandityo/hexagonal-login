package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	rh "github.com/rinosukmandityo/hexagonal-login/repositories/helper"
	"github.com/rinosukmandityo/hexagonal-login/services/logic"
)

func RegisterHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	webRepo := rh.ChooseRepo()

	registerUserHandler(r, NewUserHandler(logic.NewUserService(webRepo)))
	registerLoginHandler(r, NewLoginHandler(logic.NewLoginService(webRepo)))

	return r
}

func registerUserHandler(r *chi.Mux, handler UserHandler) {
	r.Get("/user/{id}", handler.Get)
	r.Post("/user", handler.Post)
	r.Put("/user/{id}", handler.Update)
	r.Delete("/user/{id}", handler.Delete)
}

func registerLoginHandler(r *chi.Mux, handler LoginHandler) {
	r.Post("/auth", handler.Auth)
}
