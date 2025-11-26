//nolint:errcheck
package api

import (
	"encoding/json"
	"foodjiassignment/internal/api/errors"
	"foodjiassignment/internal/api/models"
	"net/http"
)

func (a *APIInterface) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sess, err := a.tinderFoodMgr.CreateSession()
	if err != nil {
		a.apiError.FailWithMessage(w, errors.Error{
			Message:    "failed to encode response",
			StatusCode: http.StatusInternalServerError,
		})

		return
	}

	resp := models.SessionResponse{
		SessionId: sess,
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
