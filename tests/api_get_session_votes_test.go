package tests

import (
	"bytes"
	"encoding/json"
	"foodjiassignment/internal/app"
	"foodjiassignment/internal/repository/models"
	"foodjiassignment/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProductVotesPerSession(t *testing.T) {
	cfg := defaultConfig()

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(db, cfg)

	ts := httptest.NewServer(a.API.RegisterHandlers())
	defer ts.Close()

	// first create a unique session
	sessionModel := models.Session{
		ID: uuid.New(),
	}
	db.Create(&sessionModel)

	productId1 := uuid.New()
	productId2 := uuid.New()
	productId3 := uuid.New()

	bodies := []map[string]interface{}{
		{
			"productId":   productId1.String(),
			"productName": "name-1",
			"liked":       false,
		},
		{
			"productId":   productId2.String(),
			"productName": "name-2",
			"liked":       true,
		},
		{
			"productId":   productId3.String(),
			"productName": "name-3",
			"liked":       false,
		},
	}

	for _, body := range bodies {
		payload, _ := json.Marshal(body)

		upsert(ts.URL, payload, sessionModel.ID.String(), t)
	}

	// after voting for products lets get votes and check count
	votes := getProductVotesPerSession(t, ts.URL, sessionModel.ID.String())

	assert.Equal(t, len(bodies), len(votes))

	// lets now try to update one of our product with same session and recheck
	payload, _ := json.Marshal(map[string]interface{}{
		"productId":   productId2.String(),
		"productName": "name-2",
		"liked":       false,
	})
	upsert(ts.URL, payload, sessionModel.ID.String(), t)

	// we recheck the count it should still be the same, session didnt change
	votes = getProductVotesPerSession(t, ts.URL, sessionModel.ID.String())

	assert.Equal(t, len(bodies), len(votes))

	// lets create a second session id and retry for both sessions

	// first create a unique session
	sessionModel2 := models.Session{
		ID: uuid.New(),
	}
	db.Create(&sessionModel2)

	bodies2 := []map[string]interface{}{
		{
			"productId":   productId1.String(),
			"productName": "name-1",
			"liked":       false,
		},
		{
			"productId":   productId2.String(),
			"productName": "name-2",
			"liked":       true,
		},
		{
			"productId":   productId3.String(),
			"productName": "name-3",
			"liked":       false,
		},
	}

	for _, body2 := range bodies2 {
		p, _ := json.Marshal(body2)

		upsert(ts.URL, p, sessionModel2.ID.String(), t)
	}

	// after voting for products with a separate session lets get votes and check count
	votesForSession2 := getProductVotesPerSession(t, ts.URL, sessionModel2.ID.String())

	assert.Equal(t, len(bodies2), len(votesForSession2))

	// lets recheck for previous session again
	votes = getProductVotesPerSession(t, ts.URL, sessionModel.ID.String())

	assert.Equal(t, len(bodies), len(votes))
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
