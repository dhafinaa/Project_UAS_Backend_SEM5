package middleware

import (
	"errors"
	"PROJECT_UAS/helper"
	"strings"
)

type AuthContext struct {
	UserID      string
	Role        string
	Permissions []string
}

func ExtractToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid authorization format")
	}

	return parts[1], nil
}

func ValidateToken(header string) (*AuthContext, error) {
	tokenString, err := ExtractToken(header)
	if err != nil {
		return nil, err
	}

	claims, err := helper.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &AuthContext{
		UserID:      claims.ID,
		Role:        claims.Role,
		Permissions: claims.Permissions,
	}, nil
}

func RequirePermission(ctx *AuthContext, required string) error {
	if helper.HasPermission(required, ctx.Permissions) {
		return nil
	}
	return errors.New("permission denied")
}

func RequireRole(ctx *AuthContext, role string) error {
	if ctx.Role == role {
		return nil
	}
	return errors.New("role denied")
}