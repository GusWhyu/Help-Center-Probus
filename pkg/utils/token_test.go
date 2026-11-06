package utils

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateTokenRejectsNonHS256SigningMethod(t *testing.T) {
	claims := &JWTClaims{UserID: "1"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte("test-secret-that-is-long-enough-123"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	if _, err := NewTokenService("test-secret-that-is-long-enough-123", 1).ValidateToken(signedToken); err == nil {
		t.Fatal("ValidateToken() accepted a token not signed with HS256")
	}
}
