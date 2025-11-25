package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Manager struct {
	tinderFoodRepo repo
}

func NewManager(
	tinderFoodRepo repo,
) *Manager {
	return &Manager{
		tinderFoodRepo: tinderFoodRepo,
	}
}

func (a *Manager) RegisterHandlers() http.Handler {
	r := chi.NewRouter()

	r.Post("/session", a.CreateSession)

	return r
}
