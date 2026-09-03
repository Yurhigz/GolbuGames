package api_errors

import "net/http"

var (
    ErrAuth      = &sentinelAPIError{status: http.StatusUnauthorized, msg: "invalid token"}
    ErrNotFound  = &sentinelAPIError{status: http.StatusNotFound, msg: "not found"}
    ErrDuplicate = &sentinelAPIError{status: http.StatusConflict, msg: "duplicate"}
	ErrBadRequest = &sentinelAPIError{status: http.StatusBadRequest, msg: "bad request"}
	ErrForbidden = &sentinelAPIError{status: http.StatusForbidden, msg: "forbidden"}
	ErrTooManyRequests = &sentinelAPIError{status: http.StatusTooManyRequests, msg: "too many requests"}
)

type APIError interface {
	APIError() (int,string)
}

type sentinelAPIError struct {
    status int
    msg    string
}

func (e sentinelAPIError) Error() string {
    return e.msg
}

func (e sentinelAPIError) APIError() (int, string) {
    return e.status, e.msg
}



func WriteError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(APIError); ok {
		status, msg := apiErr.APIError()
		http.Error(w, msg, status)
	} else {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}