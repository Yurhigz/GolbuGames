package middleware

import (
	"fmt"
	"os"
	"time"
    "net/http"
    "context"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"strconv"
)

type CustomClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

var secretKey = []byte(os.Getenv("SecretKey"))

var refreshSecretKey = []byte(os.Getenv("RefreshSecretKey"))

func GenerateJWT(userID, username string, roles []string) (string, error) {

	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "golbugames",
			Subject:   userID,
			ID:        "golbugames_jwt_" + userID + "_" + time.Now().Format("20060102150405"),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyAndExtractClaims(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
			}
			return secretKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("erreur d'analyse du token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalide")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token expiré")
	}

	return claims, nil
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 jours
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "golbugames",
		Subject:   userID,
		ID:        "golbugames_refresh_" + userID + "_" + time.Now().Format("20060102150405"),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecretKey)
}

func VerifyRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
		}
		return refreshSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("refresh token invalide ou expiré")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expiré")
	}

	return claims, nil
}

func HasRole(claims *CustomClaims, requiredRole string) bool {
	for _, role := range claims.Roles {
		if role == requiredRole {
			return true
		}
	}
	return false
}

func JWTMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		if requiredRole != "" && !HasRole(claims, requiredRole) {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func GetUserIdFromClaim(r *http.Request) (int, error) {
    claims, ok := r.Context().Value("claims").(*CustomClaims)
    if !ok || claims == nil {
        return 0, fmt.Errorf("unauthorized: no claims found in context")
    }

    userID, err := strconv.Atoi(claims.UserID)
    if err != nil {
        return 0, fmt.Errorf("invalid user ID: %v", err)
    }

    return userID, nil
}