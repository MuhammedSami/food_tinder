package tests

import (
	"bytes"
	"encoding/json"
	"foodtinder/internal/app"
	"foodtinder/internal/repository/models"
	"foodtinder/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProductVotesPerSession(t *testing.T) {
	ResetPrometheusRegistry()

	cfg := defaultConfig()
	db := storage.NewDb(cfg.DB)
	a := app.NewApp(db, cfg)

	db.Exec("DELETE FROM product_votes")
	db.Exec("DELETE FROM sessions")

	ts := httptest.NewServer(a.API.RegisterHandlers())
	defer ts.Close()

	t.Run("create session and insert votes", func(t *testing.T) {
		session1 := models.Session{ID: uuid.New()}
		db.Create(&session1)

		// Test data
		productId1 := uuid.New()
		productId2 := uuid.New()
		productId3 := uuid.New()

		bodies := []map[string]interface{}{
			{"productId": productId1.String(), "productName": "name-1", "liked": false},
			{"productId": productId2.String(), "productName": "name-2", "liked": true},
			{"productId": productId3.String(), "productName": "name-3", "liked": false},
		}

		for _, b := range bodies {
			p, _ := json.Marshal(b)
			upsert(ts.URL, p, session1.ID.String(), t)
		}

		votes := getProductVotesPerSession(t, ts.URL, session1.ID.String())
		assert.Equal(t, len(bodies), len(votes), "expected initial vote count == 3")
	})

	t.Run("update existing product vote but keep same count", func(t *testing.T) {
		session1 := models.Session{}
		db.First(&session1) // fetch previous session

		productId2 := "" // find the productID by reading existing votes
		votes := getProductVotesPerSession(t, ts.URL, session1.ID.String())
		productId2 = votes[1]["productId"].(string)

		payload, _ := json.Marshal(map[string]interface{}{
			"productId":   productId2,
			"productName": "name-2",
			"liked":       false,
		})

		upsert(ts.URL, payload, session1.ID.String(), t)

		votesAfter := getProductVotesPerSession(t, ts.URL, session1.ID.String())
		assert.Equal(t, 3, len(votesAfter), "count must remain the same after update")
	})

	t.Run("create second session and insert votes", func(t *testing.T) {
		session2 := models.Session{ID: uuid.New()}
		db.Create(&session2)

		// Reuse product IDs for test consistency
		productId1 := uuid.New()
		productId2 := uuid.New()
		productId3 := uuid.New()

		bodies2 := []map[string]interface{}{
			{"productId": productId1.String(), "productName": "name-1", "liked": false},
			{"productId": productId2.String(), "productName": "name-2", "liked": true},
			{"productId": productId3.String(), "productName": "name-3", "liked": false},
		}

		for _, b := range bodies2 {
			p, _ := json.Marshal(b)
			upsert(ts.URL, p, session2.ID.String(), t)
		}

		votes := getProductVotesPerSession(t, ts.URL, session2.ID.String())
		assert.Equal(t, 3, len(votes), "second session must also have 3 votes")
	})

	t.Run("ensure first session unchanged after second session votes", func(t *testing.T) {
		var session1 models.Session
		db.First(&session1)

		votes := getProductVotesPerSession(t, ts.URL, session1.ID.String())
		assert.Equal(t, 3, len(votes), "first session must remain unchanged")
	})
}

func upsert(url string, payload []byte, sessionId string, t *testing.T) {
	req, _ := http.NewRequest("POST", url+"/product-votes", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", sessionId)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func getProductVotesPerSession(t *testing.T, url string, sessionId string) []map[string]interface{} {
	req, _ := http.NewRequest("GET", url+"/product-votes", nil)
	req.Header.Set("X-Session-ID", sessionId)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, _ := io.ReadAll(resp.Body)

	var votes []map[string]interface{}
	err = json.Unmarshal(responseBody, &votes)
	require.NoError(t, err)

	return votes
}
