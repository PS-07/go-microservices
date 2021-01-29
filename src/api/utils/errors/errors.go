package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

// APIError interdace
type APIError interface {
	StatusFn() int
	MessageFn() string
	ErrorFn() string
}

type apiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (e *apiError) StatusFn() int     { return e.Status }
func (e *apiError) MessageFn() string { return e.Message }
func (e *apiError) ErrorFn() string   { return e.Error }

// NewAPIError func
func NewAPIError(statusCode int, message string) APIError {
	return &apiError{
		Status:  statusCode,
		Message: message,
	}
}

// NewAPIErrFromBytes func
func NewAPIErrFromBytes(body []byte) (APIError, error) {
	var result apiError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("invalid json body")
	}
	return &result, nil
}

// NewInternalServerError func
func NewInternalServerError(message string) APIError {
	return &apiError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

// NewNotFoundError func
func NewNotFoundError(message string) APIError {
	return &apiError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

// NewBadRequestError func
func NewBadRequestError(message string) APIError {
	return &apiError{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}
