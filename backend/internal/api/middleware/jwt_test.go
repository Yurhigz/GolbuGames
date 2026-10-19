package middleware

import (
	"testing"
	"time"

	"golbugames/config"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTManagerGeneratesAndVerifiesToken(t *testing.T) {
	manager, err := NewJWTManager(config.JWTSettings{
		Secret: "0123456789abcdef0123456789abcdef",
		Issuer: "golbugames-test",
		TTL:    time.Minute,
	})
	if err != nil {
		t.Fatalf("create JWT manager: %v", err)
	}

	token, err := manager.Generate("user-7", "alice", []string{"player"})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	claims, err := manager.VerifyAndExtractClaims(token)
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if claims.UserID != "user-7" || claims.Subject != "user-7" || claims.Username != "alice" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if claims.Issuer != "golbugames-test" || !HasRole(claims, "player") {
		t.Fatalf("issuer or roles missing from claims: %+v", claims)
	}
}

func TestJWTManagerRejectsWrongIssuerAndSigningMethod(t *testing.T) {
	settings := config.JWTSettings{
		Secret: "0123456789abcdef0123456789abcdef",
		Issuer: "expected-issuer",
		TTL:    time.Minute,
	}
	manager, err := NewJWTManager(settings)
	if err != nil {
		t.Fatalf("create JWT manager: %v", err)
	}

	wrongIssuerToken := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: "user-7",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "other-issuer",
			Subject:   "user-7",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	wrongIssuer, err := wrongIssuerToken.SignedString([]byte(settings.Secret))
	if err != nil {
		t.Fatalf("sign wrong-issuer token: %v", err)
	}
	if _, err := manager.VerifyAndExtractClaims(wrongIssuer); err == nil {
		t.Fatal("expected token with wrong issuer to be rejected")
	}

	wrongMethodToken := jwt.NewWithClaims(jwt.SigningMethodHS512, CustomClaims{
		UserID: "user-7",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    settings.Issuer,
			Subject:   "user-7",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	wrongMethod, err := wrongMethodToken.SignedString([]byte(settings.Secret))
	if err != nil {
		t.Fatalf("sign wrong-method token: %v", err)
	}
	if _, err := manager.VerifyAndExtractClaims(wrongMethod); err == nil {
		t.Fatal("expected non-HS256 token to be rejected")
	}
}

func TestNewJWTManagerRejectsWeakKey(t *testing.T) {
	_, err := NewJWTManager(config.JWTSettings{
		Secret: "weak",
		Issuer: "golbugames",
		TTL:    time.Hour,
	})
	if err == nil {
		t.Fatal("expected weak signing key to be rejected")
	}
}
