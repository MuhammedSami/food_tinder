package api

import (
	"context"
	"errors"
	apiErrors "foodtinder/internal/api/errors"
	repoErrors "foodtinder/internal/repository/errors"
	chiprom "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"net/http"
)

const sessionIDKey = "sessionID"

type APIInterface struct {
	tinderFoodMgr manager
	apiError      *apiErrors.ApiError
	redis         *redis.Client
}

func NewApi(
	tinderFoodMgr manager,
	errTransformer *apiErrors.ApiError,
	redis *redis.Client,
) *APIInterface {
	return &APIInterface{
		tinderFoodMgr: tinderFoodMgr,
		apiError:      errTransformer,
		redis:         redis,
	}
}

func (a *APIInterface) RegisterHandlers() http.Handler {
	r := chi.NewRouter()

	metricMiddleware := chiprom.NewMiddleware("food_tinder")
	r.Use(metricMiddleware)

	r.Handle("/metrics", promhttp.Handler())

	r.Post("/sessions", a.CreateSession)
	r.Get("/product-scores", a.GetAverageScores)

	r.Group(func(pr chi.Router) {
		pr.Use(a.RequireSession)

		pr.Post("/product-votes", a.Upsert)
		pr.Get("/product-votes", a.GetVotesForSession)
	})

	return r
}

func (a *APIInterface) RequireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionIDStr := r.Header.Get("X-Session-ID")

		if sessionIDStr == "" {
			a.apiError.FailWithMessage(w, apiErrors.Error{
				Message:    "missing session ID",
				StatusCode: http.StatusBadRequest,
			})

			return
		}

		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			a.apiError.FailWithMessage(w, apiErrors.Error{
				Message:    "invalid session ID",
				StatusCode: http.StatusBadRequest,
			})

			return
		}

		session, err := a.tinderFoodMgr.GetSession(sessionID)
		if err != nil {
			if errors.As(err, &repoErrors.NotFound{}) {
				a.apiError.FailWithMessage(w, apiErrors.Error{
					Message:    "internal server error",
					StatusCode: http.StatusInternalServerError,
				})

				return
			}

			a.apiError.FailWithMessage(w, apiErrors.Error{
				Message:    "internal server error",
				StatusCode: http.StatusInternalServerError,
			})

			return
		}

		ctx := context.WithValue(r.Context(), sessionIDKey, session.ID.String()) //nolint:staticcheck
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
