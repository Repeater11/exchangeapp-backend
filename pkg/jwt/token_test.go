package jwt

import (
	"testing"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func TestGenerateTokenEmptySecret(t *testing.T) {
	if _, err := GenerateToken(1, "alice", "", 60); err == nil {
		t.Fatalf("expected error for empty secret")
	}
}

func TestGenerateTokenDefaultExpire(t *testing.T) {
	tokenStr, err := GenerateToken(1, "alice", "secret", 0)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}
	claims, err := ParseToken(tokenStr, "secret")
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}
	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		t.Fatalf("expected iat and exp to be set")
	}

	diff := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time)
	if diff < 59*time.Minute || diff > 61*time.Minute {
		t.Fatalf("unexpected expire duration: %v", diff)
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	tokenStr, err := GenerateToken(1, "alice", "secret", 60)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}
	if _, err := ParseToken(tokenStr, "other"); err == nil {
		t.Fatalf("expected error for wrong secret")
	}
}

func TestParseTokenExpired(t *testing.T) {
	now := time.Now().Add(-2 * time.Hour)
	claims := Claims{
		UserID:   1,
		Username: "alice",
		RegisteredClaims: jwtv5.RegisteredClaims{
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(time.Minute)),
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token failed: %v", err)
	}

	if _, err := ParseToken(tokenStr, "secret"); err == nil {
		t.Fatalf("expected error for expired token")
	}
}

func TestParseTokenEmpty(t *testing.T) {
	if _, err := ParseToken("", "secret"); err == nil {
		t.Fatalf("expected error for empty token")
	}
}
