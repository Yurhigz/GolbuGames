package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golbugames/config"
)

func TestNewHTTPHandlerCORS(t *testing.T) {
	handler := NewHTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), middlewareSettings())

	t.Run("allowed preflight", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api", nil)
		req.Header.Set("Origin", "https://game.example")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://game.example" {
			t.Errorf("expected exact allowed origin, got %q", got)
		}
	})

	t.Run("disallowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set("Origin", "https://attacker.example")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestNewHTTPHandlerLimitsRequestBody(t *testing.T) {
	settings := middlewareSettings()
	settings.Payload.HTTPMaxBytes = 4
	handler := NewHTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	}), settings)
	req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader("12345"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}
}

func TestClientLimiterRefillsTokens(t *testing.T) {
	limiter := &clientLimiter{
		buckets:    make(map[string]clientBucket),
		rate:       1,
		burst:      1,
		maxClients: 1,
		stateTTL:   time.Minute,
	}
	now := time.Now()

	if !limiter.allow("client", now) {
		t.Fatal("first request should be allowed")
	}
	if limiter.allow("client", now) {
		t.Fatal("second immediate request should be limited")
	}
	if !limiter.allow("client", now.Add(time.Second)) {
		t.Fatal("request should be allowed after one token refills")
	}
}

func middlewareSettings() config.AppSettings {
	return config.AppSettings{
		CORS: config.CORSSettings{AllowedOrigins: []string{"https://game.example"}},
		Payload: config.PayloadSettings{
			HTTPMaxBytes: 1024,
		},
		RateLimit: config.RateLimitSettings{
			RequestsPerSecond: 10,
			Burst:             20,
			MaxTrackedClients: 10,
			ClientStateTTL:    time.Minute,
		},
	}
}
