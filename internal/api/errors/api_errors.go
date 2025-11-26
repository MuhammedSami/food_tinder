package errors

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func (e *Error) ToString() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		return []byte(`{"message":"failed to encode error"}`)
	}

	return data
}

type ApiError struct{}

func NewApiError() *ApiError {
	return &ApiError{}
}

func (e *ApiError) FailWithMessage(w http.ResponseWriter, err Error) {
	w.WriteHeader(err.StatusCode)
	w.Write(err.ToString())
}
