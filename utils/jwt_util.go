package utils

import (
	"errors"
	"time"

	"whatsapp-backend/models/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct {
	Role     string `json:"role"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

type TokenUser struct {
	ID       uuid.UUID
	Role     entity.UserRole
	Email    string
	FullName string
}

func GenerateAccessToken(secret string, user TokenUser, ttl time.Duration) (string, int64, error) {
	if secret == "" {
		return "", 0, errors.New("jwt secret is empty")
	}

	now := time.Now()
	expiresAt := now.Add(ttl)
	claims := AccessClaims{
		Role:     string(user.Role),
		Email:    user.Email,
		FullName: user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        uuid.NewString(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return token, int64(ttl.Seconds()), nil
}

func ParseAccessToken(secret string, tokenString string) (*AccessClaims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is empty")
	}

	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.Subject == "" || claims.Role == "" {
		return nil, errors.New("missing required token claims")
	}
	return claims, nil
}
