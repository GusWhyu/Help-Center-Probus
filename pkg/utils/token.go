package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   string         `json:"user_id"`
	Username string         `json:"username"`
	Role     string         `json:"role"`
	Extra    map[string]any `json:"extra,omitempty"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secretKey     []byte
	lifetimeHours int
}

func NewTokenService(secretKey string, lifetimeHours int) *TokenService {
	return &TokenService{
		secretKey:     []byte(secretKey),
		lifetimeHours: lifetimeHours,
	}
}

func (s *TokenService) GenerateToken(userID, username, role string, extra map[string]any) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.lifetimeHours) * time.Hour)

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Extra:    extra,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *TokenService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid or expired token")
}
