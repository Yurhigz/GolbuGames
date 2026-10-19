package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadUsesValidatedDefaults(t *testing.T) {
	settings, err := loadFrom(testEnvironment(map[string]string{
		"DATABASE_URL": "postgres://golbu:password@localhost:5442/golbugamesdb?sslmode=disable",
		"JWT_SECRET":   "0123456789abcdef0123456789abcdef",
	}))
	if err != nil {
		t.Fatalf("load settings: %v", err)
	}

	if settings.HTTP.Address != ":3002" {
		t.Errorf("expected default address :3002, got %q", settings.HTTP.Address)
	}
	if settings.JWT.TTL != 24*time.Hour {
		t.Errorf("expected 24h JWT TTL, got %s", settings.JWT.TTL)
	}
	if len(settings.CORS.AllowedOrigins) != 2 {
		t.Errorf("expected two development origins, got %v", settings.CORS.AllowedOrigins)
	}
	if settings.Payload.HTTPMaxBytes != 1<<20 || settings.Payload.WebSocketMaxBytes != 1<<20 {
		t.Errorf("expected 1 MiB payload defaults, got HTTP=%d WebSocket=%d", settings.Payload.HTTPMaxBytes, settings.Payload.WebSocketMaxBytes)
	}
}

func TestLoadRejectsInvalidConfiguration(t *testing.T) {
	testCases := []struct {
		name       string
		overrides  map[string]string
		removeKeys []string
		wantError  string
	}{
		{name: "missing database URL", removeKeys: []string{"DATABASE_URL"}, wantError: "DATABASE_URL is required"},
		{name: "missing JWT secret", removeKeys: []string{"JWT_SECRET"}, wantError: "JWT_SECRET is required"},
		{name: "weak JWT secret", overrides: map[string]string{"JWT_SECRET": "too-short"}, wantError: "at least 32 bytes"},
		{name: "invalid database URL", overrides: map[string]string{"DATABASE_URL": "not-a-postgres-url"}, wantError: "valid PostgreSQL"},
		{name: "wildcard origin", overrides: map[string]string{"CORS_ALLOWED_ORIGINS": "*"}, wantError: "explicit origins"},
		{name: "origin with a path", overrides: map[string]string{"CORS_ALLOWED_ORIGINS": "https://example.com/app"}, wantError: "invalid origin"},
		{name: "invalid timeout", overrides: map[string]string{"HTTP_READ_TIMEOUT": "soon"}, wantError: "HTTP_READ_TIMEOUT"},
		{name: "inconsistent pool bounds", overrides: map[string]string{"DB_POOL_MIN_CONNECTIONS": "21"}, wantError: "DB_POOL_MIN_CONNECTIONS"},
		{name: "unbounded request body", overrides: map[string]string{"HTTP_MAX_BODY_BYTES": "0"}, wantError: "HTTP_MAX_BODY_BYTES"},
		{name: "non-finite rate", overrides: map[string]string{"RATE_LIMIT_REQUESTS_PER_SECOND": "NaN"}, wantError: "finite positive number"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			environment := map[string]string{
				"DATABASE_URL": "postgres://golbu:password@localhost:5442/golbugamesdb?sslmode=disable",
				"JWT_SECRET":   "0123456789abcdef0123456789abcdef",
			}
			for key, value := range tt.overrides {
				environment[key] = value
			}
			for _, key := range tt.removeKeys {
				delete(environment, key)
			}

			_, err := loadFrom(testEnvironment(environment))
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("expected error containing %q, got %v", tt.wantError, err)
			}
		})
	}
}

func testEnvironment(values map[string]string) envLookup {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
