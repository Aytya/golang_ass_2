package handler

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"sql/repo"
)

type Handler struct {
	repo *repo.Repository
}

func NewHandler(repo *repo.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.InsertUsers)
		r.Get("/", h.FindUsers)
		r.Route("/{id}", func(r chi.Router) {
			r.Put("/", h.UpdateUserById)
			r.Delete("/", h.DeleteUserById)
		})
	})

	return r
}
