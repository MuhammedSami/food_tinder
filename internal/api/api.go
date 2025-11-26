package api

import (
	"foodjiassignment/internal/api/errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type APIInterface struct {
	tinderFoodMgr manager
	apiError      *errors.ApiError
}

func NewApi(
	tinderFoodMgr manager,
	errTransformer *errors.ApiError,
) *APIInterface {
	return &APIInterface{
		tinderFoodMgr: tinderFoodMgr,
		apiError:      errTransformer,
	}
}

func (a *APIInterface) RegisterHandlers() http.Handler {
	r := chi.NewRouter()

	r.Post("/session", a.CreateSession)

	return r
}
