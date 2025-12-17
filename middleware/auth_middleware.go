package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"PROJECT_UAS/helper"
	"PROJECT_UAS/app/repository"
)

/* Extract Bearer Token */
func extractToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization format")
	}

	return parts[1], nil
}

/* AUTH REQUIRED — sesuai FR-002 langkah 1–3 */
func AuthRequired(authRepo *repository.AuthRepository, blacklist *TokenBlacklist, ) fiber.Handler {
	return func(c *fiber.Ctx) error {

		raw := c.Get("Authorization")
		tokenString, err := extractToken(raw)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		// 🔒 cek blacklist
		if blacklist.IsBlacklisted(tokenString) {
			return fiber.NewError(fiber.StatusUnauthorized, "token revoked")
		}

		claims, err := helper.ParseToken(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}

		perms, err := authRepo.GetPermissionsByRoleName(claims.Role)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed loading permissions")
		}

		// ✅ SIMPAN KE LOCALS
		c.Locals("token", tokenString)
		c.Locals("userID", claims.ID)
		c.Locals("role", claims.Role)
		c.Locals("permissions", perms)

		return c.Next()
	}
}


/* ROLE REQUIRED */
func RoleRequired(required string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		r := c.Locals("role")
		if r == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		userRole := strings.ToLower(r.(string))
		want := strings.ToLower(required)
		if userRole != want {
			return fiber.NewError(fiber.StatusForbidden, "access denied: role mismatch")
		}
		return c.Next()
	}
}

/* PERMISSION REQUIRED */
func PermissionRequired(needed string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Locals("permissions")
		if raw == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "permissions missing")
		}
		perms := raw.([]string)
		for _, p := range perms {
			if strings.EqualFold(p, needed) {
				return c.Next()
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "access denied: permission denied")
	}
}
