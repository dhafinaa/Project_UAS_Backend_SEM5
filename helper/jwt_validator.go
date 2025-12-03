package helper

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	ID          string   `json:"id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func ParseToken(tokenString string) (*AuthClaims, error) {

	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, errors.New("cannot parse token claims")
	}

	return claims, nil
}
