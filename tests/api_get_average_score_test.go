package tests

import (
	"encoding/json"
	"foodtinder/internal/app"
	"foodtinder/internal/repository/models"
	"foodtinder/internal/storage"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ResetPrometheusRegistry() {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	prometheus.DefaultGatherer = prometheus.NewRegistry()
}

func TestGetAverageScores(t *testing.T) {
	ResetPrometheusRegistry()

	cfg := defaultConfig()
	db := storage.NewDb(cfg.DB)
	a := app.NewApp(db, cfg)

	// we will insert 4 products per session and check the average for them

	// delete all products
	db.Exec("DELETE FROM product_votes")

	ts := httptest.NewServer(a.API.RegisterHandlers())
	defer ts.Close()

	product1 := uuid.New()
	product2 := uuid.New()
	product3 := uuid.New()
	product4 := uuid.New()

	session1 := models.Session{ID: uuid.New()}
	session2 := models.Session{ID: uuid.New()}
	db.Create(&session1)
	db.Create(&session2)

	makeBody := func(id uuid.UUID, name string, liked bool) map[string]interface{} {
		return map[string]interface{}{
			"productId":   id.String(),
			"productName": name,
			"liked":       liked,
		}
	}

	t.Run("insert votes for session 1", func(t *testing.T) {
		bodies := []map[string]interface{}{
			makeBody(product1, "prod-1", true),
			makeBody(product2, "prod-2", false),
			makeBody(product3, "prod-3", true),
			makeBody(product4, "prod-4", false),
		}

		for _, b := range bodies {
			p, _ := json.Marshal(b)
			upsert(ts.URL, p, session1.ID.String(), t)
		}
	})

	t.Run("insert votes for session 2", func(t *testing.T) {

		bodies := []map[string]interface{}{
			makeBody(product1, "prod-1", true),
			makeBody(product2, "prod-2", true),
			makeBody(product3, "prod-3", false),
			makeBody(product4, "prod-4", false),
		}

		for _, b := range bodies {
			p, _ := json.Marshal(b)
			upsert(ts.URL, p, session2.ID.String(), t)
		}
	})

	t.Run("verify average scores", func(t *testing.T) {

		resp, err := http.Get(ts.URL + "/product-scores")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := io.ReadAll(resp.Body)

		var results []models.ProductScore
		err = json.Unmarshal(bodyBytes, &results)
		require.NoError(t, err)

		require.Equal(t, 4, len(results)) // I exoect 4 products

		scoreMap := make(map[string]models.ProductScore)
		for _, s := range results {
			scoreMap[s.ProductID.String()] = s
		}

		p1 := scoreMap[product1.String()]
		assert.Equal(t, 2, p1.TotalVotes)
		assert.Equal(t, 2, p1.Likes)
		assert.Equal(t, 1.0, p1.AvgScore)

		p2 := scoreMap[product2.String()]
		assert.Equal(t, 2, p2.TotalVotes)
		assert.Equal(t, 1, p2.Likes)
		assert.Equal(t, 0.5, p2.AvgScore)

		p3 := scoreMap[product3.String()]
		assert.Equal(t, 2, p3.TotalVotes)
		assert.Equal(t, 1, p3.Likes)
		assert.Equal(t, 0.5, p3.AvgScore)

		p4 := scoreMap[product4.String()]
		assert.Equal(t, 2, p4.TotalVotes)
		assert.Equal(t, 0, p4.Likes)
		assert.Equal(t, 0.0, p4.AvgScore)
	})
}
