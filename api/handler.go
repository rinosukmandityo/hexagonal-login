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
	// Subrouters:
	r.Route("/user", func(r chi.Router) {
		r.Post("/", handler.Post) // POST /user
		// Subrouters:
		r.Route("/{id}", func(r chi.Router) {
			r.Use(handler.UserCtx)
			r.Get("/", handler.Get)       // GET /user/userid01
			r.Put("/", handler.Update)    // PUT /user/userid01
			r.Delete("/", handler.Delete) // DELETE /user/userid01
		})
	})
}

func registerLoginHandler(r *chi.Mux, handler LoginHandler) {
	r.Post("/auth", handler.Auth)
}
