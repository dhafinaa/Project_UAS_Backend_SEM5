package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, roleID string, permissions []string) (string, string, error) {

	secret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// ACCESS TOKEN
	claims := jwt.MapClaims{
		"id":          userID,
		"role":        roleID,
		"permissions": permissions,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// REFRESH TOKEN
	refreshClaims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}
