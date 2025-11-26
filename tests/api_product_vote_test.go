package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"foodjiassignment/internal/api/errors"
	apiModels "foodjiassignment/internal/api/models"
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
	cfg := defaultConfig()

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(db, cfg)

	ts := httptest.NewServer(a.API.RegisterHandlers())
	defer ts.Close()

	client := &http.Client{}

	var globalSessionId string

	productId := "34d7a483-e884-49f4-a2e5-d5e3469392a8"

	tests := []struct {
		name            string
		method          string
		body            map[string]interface{}
		endpoint        string
		expectedStatus  int
		runBeforeCase   func()
		extraCaseChecks func(resp *http.Response)
	}{
		{
			name:           "Try upsert without session header",
			method:         http.MethodPost,
			endpoint:       "/product-votes/upsert",
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
			name:           "Try upsert with session id",
			method:         http.MethodPost,
			endpoint:       "/product-votes/upsert",
			expectedStatus: http.StatusOK,
			body: map[string]interface{}{
				"productId":   productId,
				"productName": "name-1",
				"like":        true,
			},
			runBeforeCase: func() {
				sessionModel := models.Session{
					ID: uuid.New(),
				}
				db.Create(&sessionModel)

				globalSessionId = sessionModel.ID.String()
			},
			extraCaseChecks: func(resp *http.Response) {
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var result apiModels.UpsertProductVoteResponse

				err = json.Unmarshal(bodyBytes, &result)
				if err != nil {
					fmt.Println("failed to unmarshal response:", err)
					return
				}

				assert.Equal(t, "vote saved for product", result.Message)
				assert.Equal(t, productId, result.ProductId)
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

			if globalSessionId != "" {
				req.Header.Set("X-Session-ID", globalSessionId)
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
