package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"foodtinder/internal/api/errors"
	apiModels "foodtinder/internal/api/models"
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

func convertResponseIntoErrorResponse(t *testing.T, resp *http.Response) errors.Error {
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result errors.Error

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		fmt.Println("failed to unmarshal response:", err)
		return result
	}

	return result
}

func TestProductVoteAPI(t *testing.T) {
	ResetPrometheusRegistry()

	cfg := defaultConfig()

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(db, cfg)

	ts := httptest.NewServer(a.API.RegisterHandlers())
	defer ts.Close()

	client := &http.Client{}

	// get a session id
	sessionModel := models.Session{
		ID: uuid.New(),
	}
	db.Create(&sessionModel)

	productId := uuid.New()
	productId2 := uuid.New()

	tests := []struct {
		name            string
		method          string
		body            map[string]interface{}
		endpoint        string
		expectedStatus  int
		sessionId       string
		runBeforeCase   func()
		extraCaseChecks func(resp *http.Response)
	}{
		{
			name:           "Try upsert without session header",
			method:         http.MethodPost,
			endpoint:       "/product-votes",
			expectedStatus: http.StatusBadRequest,
			extraCaseChecks: func(resp *http.Response) {
				result := convertResponseIntoErrorResponse(t, resp)
				assert.Equal(t, "missing session ID", result.Message)
			},
		},
		{
			name:           "Try get votes per session without session header",
			method:         http.MethodGet,
			endpoint:       "/product-votes",
			expectedStatus: http.StatusBadRequest,
			extraCaseChecks: func(resp *http.Response) {
				result := convertResponseIntoErrorResponse(t, resp)
				assert.Equal(t, "missing session ID", result.Message)
			},
		},
		{
			name:           "user can vote for a product for the first time",
			method:         http.MethodPost,
			endpoint:       "/product-votes",
			sessionId:      sessionModel.ID.String(),
			expectedStatus: http.StatusOK,
			body: map[string]interface{}{
				"productId":   productId.String(),
				"productName": "name-1",
				"liked":       true,
			},
			extraCaseChecks: func(resp *http.Response) {
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var result apiModels.UpsertProductVoteResponse

				err = json.Unmarshal(bodyBytes, &result)
				require.NoError(t, err)

				assert.Equal(t, "vote saved for product", result.Message)
				assert.Equal(t, productId.String(), result.ProductId)
			},
		},
		{
			name:           "user is able to vote for same product",
			method:         http.MethodPost,
			sessionId:      sessionModel.ID.String(),
			endpoint:       "/product-votes",
			expectedStatus: http.StatusOK,
			body: map[string]interface{}{
				"productId":   productId2.String(),
				"productName": "name-1",
				"liked":       false,
			},
			runBeforeCase: func() {
				name := "NAME-1"
				err := db.Create(&models.ProductVote{
					ProductID:   productId2,
					SessionID:   sessionModel.ID,
					ProductName: &name,
					Liked:       true,
				}).Error
				require.NoError(t, err)
			},
			extraCaseChecks: func(resp *http.Response) {
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var result apiModels.UpsertProductVoteResponse

				err = json.Unmarshal(bodyBytes, &result)
				require.NoError(t, err)

				assert.Equal(t, "vote saved for product", result.Message)
				assert.Equal(t, productId2.String(), result.ProductId)

				var model *models.ProductVote

				err = db.First(&model, "session_id = ? and product_id = ?", sessionModel.ID.String(), productId2.String()).Error
				require.NoError(t, err)
				assert.Equal(t, productId2, model.ProductID)
				assert.Equal(t, sessionModel.ID, model.SessionID)
				assert.Equal(t, false, model.Liked)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			if tt.runBeforeCase != nil {
				tt.runBeforeCase()
			}

			var req *http.Request

			switch tt.method {
			case http.MethodGet:
				req, err = http.NewRequest(http.MethodGet, ts.URL+tt.endpoint, nil)

			case http.MethodPost:
				jsonBody, _ := json.Marshal(tt.body)

				req, err = http.NewRequest(http.MethodPost, ts.URL+tt.endpoint, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			}

			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			if tt.sessionId != "" {
				req.Header.Set("X-Session-ID", tt.sessionId)
			}

			resp, err = client.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close() //nolint:errcheck

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.extraCaseChecks != nil {
				tt.extraCaseChecks(resp)
			}
		})
	}
}
