//nolint:errcheck
package api

import (
	"encoding/json"
	"foodjiassignment/internal/api/models"
	"net/http"
)

func (a *Manager) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := models.SessionResponse{
		SessionId: "",
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"failed to encode response"}`))

		return
	}
}
