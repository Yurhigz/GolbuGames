package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddTournamentRejectsMalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/tournaments", strings.NewReader(`{"name":`))
	rec := httptest.NewRecorder()

	AddTournament(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected JSON content type, got %q", contentType)
	}
	if !strings.Contains(rec.Body.String(), `"code":"BAD_REQUEST"`) {
		t.Errorf("expected BAD_REQUEST error body, got %q", rec.Body.String())
	}
}
