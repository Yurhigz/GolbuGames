package config

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppSettings struct {
	HTTP      HTTPSettings
	Database  DatabaseSettings
	JWT       JWTSettings
	CORS      CORSSettings
	Payload   PayloadSettings
	RateLimit RateLimitSettings
}

type HTTPSettings struct {
	Address           string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	MaxHeaderBytes    int
}

type DatabaseSettings struct {
	URL                  string
	MinConnections       int32
	MaxConnections       int32
	MaxConnectionAge     time.Duration
	MaxConnectionIdle    time.Duration
	ConnectTimeout       time.Duration
	QueryTimeout         time.Duration
	BulkOperationTimeout time.Duration
}

type JWTSettings struct {
	Secret string
	Issuer string
	TTL    time.Duration
}

type CORSSettings struct {
	AllowedOrigins []string
}

type PayloadSettings struct {
	HTTPMaxBytes      int64
	WebSocketMaxBytes uint64
}

type RateLimitSettings struct {
	RequestsPerSecond float64
	Burst             int
	MaxTrackedClients int
	ClientStateTTL    time.Duration
}

type envLookup func(string) (string, bool)

func Load() (AppSettings, error) {
	return loadFrom(os.LookupEnv)
}

func loadFrom(lookup envLookup) (AppSettings, error) {
	var settings AppSettings
	var err error

	settings.HTTP.Address = stringValue(lookup, "HTTP_ADDR", ":3002")
	if settings.HTTP.Address == "" {
		return settings, fmt.Errorf("HTTP_ADDR must not be empty")
	}
	if settings.HTTP.ReadTimeout, err = durationValue(lookup, "HTTP_READ_TIMEOUT", "30s"); err != nil {
		return settings, err
	}
	if settings.HTTP.ReadHeaderTimeout, err = durationValue(lookup, "HTTP_READ_HEADER_TIMEOUT", "5s"); err != nil {
		return settings, err
	}
	if settings.HTTP.WriteTimeout, err = durationValue(lookup, "HTTP_WRITE_TIMEOUT", "30s"); err != nil {
		return settings, err
	}
	if settings.HTTP.IdleTimeout, err = durationValue(lookup, "HTTP_IDLE_TIMEOUT", "60s"); err != nil {
		return settings, err
	}
	if settings.HTTP.ShutdownTimeout, err = durationValue(lookup, "HTTP_SHUTDOWN_TIMEOUT", "10s"); err != nil {
		return settings, err
	}
	if settings.HTTP.MaxHeaderBytes, err = intValue(lookup, "HTTP_MAX_HEADER_BYTES", 1<<20); err != nil {
		return settings, err
	}

	settings.Database.URL = strings.TrimSpace(stringValue(lookup, "DATABASE_URL", ""))
	if settings.Database.URL == "" {
		return settings, fmt.Errorf("DATABASE_URL is required")
	}
	if _, err = pgxpool.ParseConfig(settings.Database.URL); err != nil {
		return settings, fmt.Errorf("DATABASE_URL is not a valid PostgreSQL connection string")
	}
	if settings.Database.MinConnections, err = int32Value(lookup, "DB_POOL_MIN_CONNECTIONS", 5); err != nil {
		return settings, err
	}
	if settings.Database.MaxConnections, err = int32Value(lookup, "DB_POOL_MAX_CONNECTIONS", 20); err != nil {
		return settings, err
	}
	if settings.Database.MaxConnectionAge, err = durationValue(lookup, "DB_POOL_MAX_CONNECTION_AGE", "1h"); err != nil {
		return settings, err
	}
	if settings.Database.MaxConnectionIdle, err = durationValue(lookup, "DB_POOL_MAX_CONNECTION_IDLE", "10m"); err != nil {
		return settings, err
	}
	if settings.Database.ConnectTimeout, err = durationValue(lookup, "DB_CONNECT_TIMEOUT", "5s"); err != nil {
		return settings, err
	}
	if settings.Database.QueryTimeout, err = durationValue(lookup, "DB_QUERY_TIMEOUT", "2s"); err != nil {
		return settings, err
	}
	if settings.Database.BulkOperationTimeout, err = durationValue(lookup, "DB_BULK_OPERATION_TIMEOUT", "10s"); err != nil {
		return settings, err
	}

	settings.JWT.Secret = strings.TrimSpace(stringValue(lookup, "JWT_SECRET", ""))
	if settings.JWT.Secret == "" {
		return settings, fmt.Errorf("JWT_SECRET is required")
	}
	settings.JWT.Issuer = strings.TrimSpace(stringValue(lookup, "JWT_ISSUER", "golbugames"))
	if settings.JWT.Issuer == "" {
		return settings, fmt.Errorf("JWT_ISSUER must not be empty")
	}
	if settings.JWT.TTL, err = durationValue(lookup, "JWT_TTL", "24h"); err != nil {
		return settings, err
	}

	origins := stringValue(lookup, "CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000")
	if settings.CORS.AllowedOrigins, err = parseOrigins(origins); err != nil {
		return settings, err
	}

	if settings.Payload.HTTPMaxBytes, err = int64Value(lookup, "HTTP_MAX_BODY_BYTES", 1<<20); err != nil {
		return settings, err
	}
	if settings.Payload.WebSocketMaxBytes, err = uint64Value(lookup, "WEBSOCKET_MAX_PAYLOAD_BYTES", 1<<20); err != nil {
		return settings, err
	}

	if settings.RateLimit.RequestsPerSecond, err = floatValue(lookup, "RATE_LIMIT_REQUESTS_PER_SECOND", 10); err != nil {
		return settings, err
	}
	if settings.RateLimit.Burst, err = intValue(lookup, "RATE_LIMIT_BURST", 20); err != nil {
		return settings, err
	}
	if settings.RateLimit.MaxTrackedClients, err = intValue(lookup, "RATE_LIMIT_MAX_TRACKED_CLIENTS", 10000); err != nil {
		return settings, err
	}
	if settings.RateLimit.ClientStateTTL, err = durationValue(lookup, "RATE_LIMIT_CLIENT_STATE_TTL", "2m"); err != nil {
		return settings, err
	}

	if err := settings.Validate(); err != nil {
		return settings, err
	}
	return settings, nil
}

func (s AppSettings) Validate() error {
	if len([]byte(s.JWT.Secret)) < 32 {
		return fmt.Errorf("JWT_SECRET must contain at least 32 bytes for HS256")
	}
	if s.JWT.TTL <= 0 || s.JWT.TTL > 30*24*time.Hour {
		return fmt.Errorf("JWT_TTL must be between 0 and 720h")
	}
	if s.Database.MinConnections < 0 || s.Database.MaxConnections < 1 || s.Database.MinConnections > s.Database.MaxConnections {
		return fmt.Errorf("DB_POOL_MIN_CONNECTIONS must be between 0 and DB_POOL_MAX_CONNECTIONS")
	}
	if s.HTTP.MaxHeaderBytes < 1024 || s.HTTP.MaxHeaderBytes > 16<<20 {
		return fmt.Errorf("HTTP_MAX_HEADER_BYTES must be between 1024 and 16777216")
	}
	if s.Payload.HTTPMaxBytes < 1 || s.Payload.HTTPMaxBytes > 32<<20 {
		return fmt.Errorf("HTTP_MAX_BODY_BYTES must be between 1 and 33554432")
	}
	if s.Payload.WebSocketMaxBytes < 125 || s.Payload.WebSocketMaxBytes > 16<<20 {
		return fmt.Errorf("WEBSOCKET_MAX_PAYLOAD_BYTES must be between 125 and 16777216")
	}
	if math.IsNaN(s.RateLimit.RequestsPerSecond) || math.IsInf(s.RateLimit.RequestsPerSecond, 0) || s.RateLimit.RequestsPerSecond <= 0 {
		return fmt.Errorf("RATE_LIMIT_REQUESTS_PER_SECOND must be a finite positive number")
	}
	if s.RateLimit.Burst < 1 || s.RateLimit.MaxTrackedClients < 1 || s.RateLimit.ClientStateTTL <= 0 {
		return fmt.Errorf("rate limit burst, tracked-client limit and state TTL must be positive")
	}
	for name, value := range map[string]time.Duration{
		"HTTP_READ_TIMEOUT":           s.HTTP.ReadTimeout,
		"HTTP_READ_HEADER_TIMEOUT":    s.HTTP.ReadHeaderTimeout,
		"HTTP_WRITE_TIMEOUT":          s.HTTP.WriteTimeout,
		"HTTP_IDLE_TIMEOUT":           s.HTTP.IdleTimeout,
		"HTTP_SHUTDOWN_TIMEOUT":       s.HTTP.ShutdownTimeout,
		"DB_POOL_MAX_CONNECTION_AGE":  s.Database.MaxConnectionAge,
		"DB_POOL_MAX_CONNECTION_IDLE": s.Database.MaxConnectionIdle,
		"DB_CONNECT_TIMEOUT":          s.Database.ConnectTimeout,
		"DB_QUERY_TIMEOUT":            s.Database.QueryTimeout,
		"DB_BULK_OPERATION_TIMEOUT":   s.Database.BulkOperationTimeout,
		"RATE_LIMIT_CLIENT_STATE_TTL": s.RateLimit.ClientStateTTL,
	} {
		if value <= 0 {
			return fmt.Errorf("%s must be positive", name)
		}
	}
	return nil
}

func parseOrigins(value string) ([]string, error) {
	var origins []string
	seen := make(map[string]struct{})
	for _, item := range strings.Split(value, ",") {
		origin := strings.TrimSpace(item)
		if origin == "" {
			continue
		}
		if origin == "*" {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS must list explicit origins, not *")
		}
		parsed, err := url.Parse(origin)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS contains an invalid origin %q", origin)
		}
		if _, exists := seen[origin]; !exists {
			seen[origin] = struct{}{}
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS must contain at least one origin")
	}
	return origins, nil
}

func stringValue(lookup envLookup, name, fallback string) string {
	if value, ok := lookup(name); ok {
		return value
	}
	return fallback
}

func durationValue(lookup envLookup, name, fallback string) (time.Duration, error) {
	value := stringValue(lookup, name, fallback)
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a duration such as 5s or 1m", name)
	}
	return duration, nil
}

func intValue(lookup envLookup, name string, fallback int) (int, error) {
	value := stringValue(lookup, name, strconv.Itoa(fallback))
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", name)
	}
	return parsed, nil
}

func int32Value(lookup envLookup, name string, fallback int32) (int32, error) {
	value := stringValue(lookup, name, strconv.FormatInt(int64(fallback), 10))
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s must be a 32-bit integer", name)
	}
	return int32(parsed), nil
}

func int64Value(lookup envLookup, name string, fallback int64) (int64, error) {
	value := stringValue(lookup, name, strconv.FormatInt(fallback, 10))
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a 64-bit integer", name)
	}
	return parsed, nil
}

func uint64Value(lookup envLookup, name string, fallback uint64) (uint64, error) {
	value := stringValue(lookup, name, strconv.FormatUint(fallback, 10))
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return parsed, nil
}

func floatValue(lookup envLookup, name string, fallback float64) (float64, error) {
	value := stringValue(lookup, name, strconv.FormatFloat(fallback, 'f', -1, 64))
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", name)
	}
	return parsed, nil
}
