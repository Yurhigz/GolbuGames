package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInvalidUserStats(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/user_stats/invalid_id", nil)
	rec := httptest.NewRecorder()
	GetUserStats(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}

}
