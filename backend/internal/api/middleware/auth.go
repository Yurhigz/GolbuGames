package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type claimsContextKey struct{}

func RequireJWT(manager *JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if manager == nil {
			http.Error(w, "authentication is unavailable", http.StatusInternalServerError)
			return
		}

		header := strings.TrimSpace(r.Header.Get("Authorization"))
		scheme, tokenString, found := strings.Cut(header, " ")
		if !found || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(tokenString) == "" {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		claims, err := manager.VerifyAndExtractClaims(strings.TrimSpace(tokenString))
		if err != nil {
			http.Error(w, "invalid or expired access token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), claimsContextKey{}, claims)))
	})
}

func ClaimsFromContext(ctx context.Context) (*CustomClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*CustomClaims)
	return claims, ok && claims != nil
}

func ClaimsFromWebSocketProtocols(r *http.Request, manager *JWTManager) (*CustomClaims, error) {
	if manager == nil {
		return nil, errors.New("JWT manager is unavailable")
	}

	for _, protocol := range strings.Split(r.Header.Get("Sec-WebSocket-Protocol"), ",") {
		protocol = strings.TrimSpace(protocol)
		if tokenString, found := strings.CutPrefix(protocol, "auth."); found {
			if tokenString == "" {
				return nil, errors.New("empty WebSocket access token")
			}
			return manager.VerifyAndExtractClaims(tokenString)
		}
	}

	return nil, errors.New("WebSocket access token is missing")
}
