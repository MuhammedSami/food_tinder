//nolint:errcheck
package api

import (
	"encoding/json"
	"fmt"
	"foodjiassignment/internal/api/errors"
	"foodjiassignment/internal/api/models"
	"github.com/google/uuid"
	"net/http"
)

func (a *APIInterface) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	sessionId, err := a.tinderFoodMgr.CreateSession(ctx)
	if err != nil {
		a.apiError.FailWithMessage(w, errors.Error{
			Message:    "failed to encode response",
			StatusCode: http.StatusInternalServerError,
		})

		return
	}

	resp := models.SessionResponse{
		SessionId: sessionId,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		a.apiError.FailWithMessage(w, errors.Error{
			Message:    "failed to encode response",
			StatusCode: http.StatusInternalServerError,
		})

		return
	}
}

func (a *APIInterface) Upsert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	sessionIDVal := ctx.Value("sessionID")

	var payload models.UpsertProductVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		a.apiError.FailWithMessage(w, errors.Error{
			Message:    "invalid JSON body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if payload.ProductId == "" {
		a.apiError.FailWithMessage(w, errors.Error{
			Message:    "productId is required",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	payload.SessionId = sessionIDVal.(string)

	err := a.tinderFoodMgr.UpsertVote(ctx, payload)
	if err != nil {
		a.apiError.FailWithMessage(w, errors.Error{
			Message:    fmt.Sprintf("failed to upsert vote, err:%s", err),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	response := models.UpsertProductVoteResponse{
		ProductId: payload.ProductId,
		Message:   "vote saved for product",
	}

	respBytes, _ := json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func (a *APIInterface) GetVotesForSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	sessionIDStr := ctx.Value("sessionID")

	sessionID, err := uuid.Parse(sessionIDStr.(string))
	if err != nil {
		a.apiError.FailWithMessage(w, errors.Error{
			Message:    fmt.Sprintf("invalid session ID err: %s", err),
			StatusCode: http.StatusBadRequest,
		})

		return
	}

	votes, err := a.tinderFoodMgr.GetVotesBySession(ctx, sessionID)
	if err != nil {
		a.apiError.FailWithMessage(w, errors.Error{
			Message:    fmt.Sprintf("failed to get product votes for sesion:%s err: %s", sessionID.String(), err),
			StatusCode: http.StatusBadRequest,
		})

		return
	}

	var response []models.GetSessionProductVotesResponse

	for _, v := range votes {
		response = append(response, models.GetSessionProductVotesResponse{
			ProductId:   v.ProductID.String(),
			ProductName: v.ProductName,
			Liked:       v.Liked,
			CreatedAt:   v.CreatedAt.String(),
		})
	}

	respBytes, _ := json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func (a *APIInterface) GetAverageScores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	avgScores, err := a.tinderFoodMgr.GetAverageScores()
	if err != nil {
		a.apiError.FailWithMessage(w, errors.Error{
			Message:    fmt.Sprintf("failed to get average score per prodcut err: %s", err),
			StatusCode: http.StatusBadRequest,
		})

		return
	}

	respBytes, _ := json.Marshal(avgScores)

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
