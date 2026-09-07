package api_errors

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{name: "Authentication Error", err: ErrAuth, wantStatus: 401, wantMsg: "invalid token"},
		{name: "Not Found Error", err: ErrNotFound, wantStatus: 404, wantMsg: "not found"},
		{name: "Duplicate Error", err: ErrDuplicate, wantStatus: 409, wantMsg: "duplicate"},
		{name: "Bad Request Error", err: ErrBadRequest, wantStatus: 400, wantMsg: "bad request"},
		{name: "Forbidden Error", err: ErrForbidden, wantStatus: 403, wantMsg: "forbidden"},
		{name: "Too Many Requests Error", err: ErrTooManyRequests, wantStatus: 429, wantMsg: "too many requests"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			var resp MessageResponse
			WriteError(w, tt.err)

			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if resp.Message != tt.wantMsg {
				t.Errorf("expected message %q, got %q", tt.wantMsg, resp.Message)
			}
		})
	}
}

func TestWriteErrorUnknownError(t *testing.T) {

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{name: "Unknown Error", err: errors.New("unknown error"), wantStatus: 500, wantMsg: "internal server error"},
		{name: "Non-API Error", err: errors.New("surprise error"), wantStatus: 500, wantMsg: "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp MessageResponse
			w := httptest.NewRecorder()

			WriteError(w, tt.err)
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if resp.Message != tt.wantMsg {
				t.Errorf("expected message %q, got %q", tt.wantMsg, resp.Message)
			}

		})
	}
}
