package helper

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var Secret = []byte("SECRET_KEY")

func GenerateTokens(userID string, role string, permissions []string) (string, string) {

	claims := jwt.MapClaims{
		"id":          userID,
		"role":        role,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	access, _ := accessToken.SignedString(Secret)
	refresh, _ := refreshToken.SignedString(Secret)

	return access, refresh
}