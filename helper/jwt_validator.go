package helper

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// AuthClaims minimal (tidak menyertakan permissions)
type AuthClaims struct {
    ID          string   `json:"id"`
    Role        string   `json:"role"`
    Permissions []string `json:"permissions"`
    jwt.RegisteredClaims
}


// ParseToken parse token dan kembalikan claims (id + role)
func ParseToken(tokenString string) (*AuthClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, errors.New("cannot parse token claims")
	}

	return claims, nil
}