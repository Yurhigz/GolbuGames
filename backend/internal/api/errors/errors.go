package api_errors

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrAuth              = &sentinelAPIError{status: http.StatusUnauthorized, msg: "invalid token", code: "INVALID_TOKEN"}
	ErrNotFound          = &sentinelAPIError{status: http.StatusNotFound, msg: "not found", code: "NOT_FOUND"}
	ErrDuplicate         = &sentinelAPIError{status: http.StatusConflict, msg: "duplicate", code: "DUPLICATE"}
	ErrBadRequest        = &sentinelAPIError{status: http.StatusBadRequest, msg: "bad request", code: "BAD_REQUEST"}
	ErrForbidden         = &sentinelAPIError{status: http.StatusForbidden, msg: "forbidden", code: "FORBIDDEN"}
	ErrTooManyRequests   = &sentinelAPIError{status: http.StatusTooManyRequests, msg: "too many requests", code: "TOO_MANY_REQUESTS"}
	ErrInvalidUserID     = &sentinelAPIError{status: http.StatusBadRequest, msg: "invalid user id", code: "INVALID_USER_ID"}
	ErrInvalidFriendID   = &sentinelAPIError{status: http.StatusBadRequest, msg: "invalid friend id", code: "INVALID_FRIEND_ID"}
	ErrFriendNotFound    = &sentinelAPIError{status: http.StatusBadRequest, msg: "friend not found", code: "FRIEND_NOT_FOUND"}
	ErrInvalidDifficulty = &sentinelAPIError{status: http.StatusBadRequest, msg: "invalid difficulty level", code: "INVALID_DIFFICULTY_LEVEL"}
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type APIErrorDetails struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type APIError interface {
	APIError() APIErrorDetails
}

type sentinelAPIError struct {
	status int
	msg    string
	code   string
}

func (e sentinelAPIError) Error() string {
	return e.msg
}

func (e sentinelAPIError) APIError() APIErrorDetails {
	return APIErrorDetails{
		Status:  e.status,
		Message: e.msg,
		Code:    e.code,
	}
}

func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := "internal server error"
	code := "INTERNAL_SERVER_ERROR"

	var APIError APIError

	if errors.As(err, &APIError) {
		apiErrDetails := APIError.APIError()
		status = apiErrDetails.Status
		msg = apiErrDetails.Message
		code = apiErrDetails.Code
	}

	response := ErrorResponse{Code: code, Message: msg}

	body, encodeErr := json.Marshal(response)
	if encodeErr != nil {
		slog.Error("failed to encode error response", "error", encodeErr)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, writeErr := w.Write(body); writeErr != nil {
		slog.Error("failed to write error response", "error", writeErr)
	}

}
