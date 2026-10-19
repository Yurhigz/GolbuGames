package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	api_errors "golbugames/internal/api/errors"
)

func TestAddGridRejectsInvalidRequests(t *testing.T) {
	testCases := []struct {
		name          string
		body          string
		wantStatus    int
		wantErrorCode string
	}{
		{
			name:          "unsupported difficulty",
			body:          `{"difficulty":"impossible"}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: "INVALID_DIFFICULTY_LEVEL",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/sudoku", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			AddGrid(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response api_errors.ErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if response.Code != tt.wantErrorCode {
				t.Errorf("expected error code %q, got %q", tt.wantErrorCode, response.Code)
			}
		})
	}
}

func TestGetGridRejectsUnsupportedDifficulty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/sudoku?difficulty=impossible", nil)
	rec := httptest.NewRecorder()

	GetGrid(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
