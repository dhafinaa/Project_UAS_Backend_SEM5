package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken membuat access + refresh token.
// Token access hanya menyimpan user id dan role (nama role).
func GenerateToken(userID string, roleName string, permissions []string) (string, string, error) {
    secret := os.Getenv("JWT_SECRET")
    refreshSecret := os.Getenv("REFRESH_SECRET")

    claims := jwt.MapClaims{
        "id":          userID,
        "role":        roleName,
        "permissions": permissions,   // <--- WAJIB ADA
        "exp":         time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    accessStr, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", "", err
    }

    refreshClaims := jwt.MapClaims{
        "id":  userID,
        "exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshStr, err := refreshToken.SignedString([]byte(refreshSecret))

    return accessStr, refreshStr, err
}
