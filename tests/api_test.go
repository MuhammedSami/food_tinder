package tests

import (
	"encoding/json"
	"fmt"
	"foodjiassignment/config"
	"foodjiassignment/internal/app"
	"foodjiassignment/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIEndpoints(t *testing.T) {
	cfg := config.Config{
		Api: config.API{
			Port: 8080,
		},
		DB: config.DBConn{
			Host:     "localhost",
			User:     "tinder_user",
			Password: "tinder_pass", // load the pass in a more secure way
			DBName:   "tinder_db",
			Port:     5432,
		},
	}

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(db, &cfg)

	ts := httptest.NewServer(a.API.RegisterHandlers())
	defer ts.Close()

	tests := []struct {
		name            string
		method          string
		endpoint        string
		expectedStatus  int
		extraCaseChecks func(resp *http.Response)
	}{
		{
			name:           "Get session id",
			method:         http.MethodPost,
			endpoint:       "/session",
			expectedStatus: http.StatusOK,
			extraCaseChecks: func(resp *http.Response) {
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var result struct {
					SessionId string `json:"sessionId"`
				}

				err = json.Unmarshal(bodyBytes, &result)
				if err != nil {
					fmt.Println("failed to unmarshal response:", err)
					return
				}

				assert.NotNil(t, result.SessionId)

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			switch tt.method {
			case http.MethodGet:
				resp, err = http.Get(ts.URL + tt.endpoint)
			case http.MethodPost:
				resp, err = http.Post(ts.URL+tt.endpoint, "application/json", nil)
			}

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
