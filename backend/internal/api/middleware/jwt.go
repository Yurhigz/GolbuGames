package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golbugames/config"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTManager(settings config.JWTSettings) (*JWTManager, error) {
	if len([]byte(settings.Secret)) < 32 {
		return nil, fmt.Errorf("JWT secret must contain at least 32 bytes")
	}
	if settings.Issuer == "" {
		return nil, fmt.Errorf("JWT issuer must not be empty")
	}
	if settings.TTL <= 0 {
		return nil, fmt.Errorf("JWT TTL must be positive")
	}
	return &JWTManager{
		secret: []byte(settings.Secret),
		issuer: settings.Issuer,
		ttl:    settings.TTL,
	}, nil
}

func (m *JWTManager) Generate(userID, username string, roles []string) (string, error) {
	now := time.Now()
	identifier := make([]byte, 16)
	if _, err := rand.Read(identifier); err != nil {
		return "", fmt.Errorf("generate JWT identifier: %w", err)
	}

	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Subject:   userID,
			ID:        hex.EncodeToString(identifier),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *JWTManager) VerifyAndExtractClaims(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, fmt.Errorf("erreur d'analyse du token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalide")
	}

	if claims.UserID == "" || claims.Subject != claims.UserID {
		return nil, fmt.Errorf("token subject does not match user ID")
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

// func VerifyToken(tokenString string) error {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	if !token.Valid {
// 		return fmt.Errorf("invalid token")
// 	}

// 	return nil
// }

// https://medium.com/@cheickzida/golang-implementing-jwt-token-authentication-bba9bfd84d60
// https://www.sohamkamani.com/golang/jwt-authentication/
