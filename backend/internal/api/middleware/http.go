package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golbugames/config"
)

type clientBucket struct {
	tokens    float64
	updatedAt time.Time
	seenAt    time.Time
}

type clientLimiter struct {
	mu          sync.Mutex
	buckets     map[string]clientBucket
	rate        float64
	burst       float64
	maxClients  int
	stateTTL    time.Duration
	cleanupTime time.Time
}

func NewHTTPHandler(next http.Handler, settings config.AppSettings) http.Handler {
	limiter := &clientLimiter{
		buckets:    make(map[string]clientBucket),
		rate:       settings.RateLimit.RequestsPerSecond,
		burst:      float64(settings.RateLimit.Burst),
		maxClients: settings.RateLimit.MaxTrackedClients,
		stateTTL:   settings.RateLimit.ClientStateTTL,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if !containsOrigin(settings.CORS.AllowedOrigins, origin) {
				http.Error(w, "origin not allowed", http.StatusForbidden)
				return
			}
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Max-Age", "600")
		}

		if r.Method == http.MethodOptions {
			if origin == "" {
				http.Error(w, "origin header required for preflight", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if !limiter.allow(clientIP(r), time.Now()) {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/ws/") && r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, settings.Payload.HTTPMaxBytes)
		}
		next.ServeHTTP(w, r)
	})
}

func containsOrigin(origins []string, origin string) bool {
	for _, allowed := range origins {
		if origin == allowed {
			return true
		}
	}
	return false
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func (l *clientLimiter) allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if now.Sub(l.cleanupTime) >= l.stateTTL {
		for client, bucket := range l.buckets {
			if now.Sub(bucket.seenAt) >= l.stateTTL {
				delete(l.buckets, client)
			}
		}
		l.cleanupTime = now
	}

	bucket, exists := l.buckets[key]
	if !exists {
		if len(l.buckets) >= l.maxClients {
			return false
		}
		bucket = clientBucket{tokens: l.burst, updatedAt: now}
	}

	elapsed := now.Sub(bucket.updatedAt).Seconds()
	bucket.tokens += elapsed * l.rate
	if bucket.tokens > l.burst {
		bucket.tokens = l.burst
	}
	bucket.updatedAt = now
	if bucket.tokens < 1 {
		l.buckets[key] = bucket
		return false
	}
	bucket.tokens--
	bucket.seenAt = now
	l.buckets[key] = bucket
	return true
}
