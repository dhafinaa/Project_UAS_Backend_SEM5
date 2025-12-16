package service

import (

	"PROJECT_UAS/middleware"
	"PROJECT_UAS/app/model"
	"PROJECT_UAS/app/repository"
	"PROJECT_UAS/helper"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	Repo *repository.AuthRepository
	blacklist *middleware.TokenBlacklist
}

func NewAuthService(repo *repository.AuthRepository, blacklist *middleware.TokenBlacklist) *AuthService {
	return &AuthService{Repo: repo, blacklist: blacklist,}
}

//
// ------------------------------------------------------
// LOGIN (FR-001) — Lengkap dengan load permissions (FR-002 Step 3)
// ------------------------------------------------------
func (s *AuthService) Login(c *fiber.Ctx) error {

	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	// 1. Ambil user berdasarkan username/email
	user, roleName, err := s.Repo.FindByLogin(req.Login)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid username/email or password")
	}

	// 2. User inactive?
	if !user.Is_active {
		return fiber.NewError(fiber.StatusUnauthorized, "user inactive")
	}

	// 3. Cek password
	if !helper.CheckPasswordHash(req.Password, user.Password_hash) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid username/email or password")
	}

	// 4. Load permissions berdasarkan ROLE NAME (SRS FR-002 step 3)
	perms, err := s.Repo.GetPermissionsByRoleName(roleName)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed loading permissions")
	}

	// 5. Generate token (SRS: token memuat id + role + permissions)
	access, refresh, err := helper.GenerateToken(user.ID, roleName, perms)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed generating token")
	}

	// 6. Return JSON sesuai format SRS
	return c.JSON(model.LoginResponse{
		Token:       access,
		Refresh:     refresh,
		User:        *user,
		Permissions: perms,
	})
}

//
// ------------------------------------------------------
// PROFILE (FR-002 – setelah middleware load user + permissions)
// ------------------------------------------------------
func (s *AuthService) Profile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"userID":      c.Locals("userID"),
		"role":        c.Locals("role"),
		"permissions": c.Locals("permissions"),
	})
}

//
// ------------------------------------------------------
// REFRESH TOKEN (versi sederhana sesuai SRS)
// ------------------------------------------------------
func (s *AuthService) RefreshToken(c *fiber.Ctx) error {

	var req struct {
		Refresh string `json:"refresh"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	claims, err := helper.ParseToken(req.Refresh)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
	}

	// Ambil user + role (tanpa permissions)
	user, roleName, _, err := s.Repo.GetUserRoleByID(claims.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "user not found")
	}

	// Buat token baru (permissions tidak wajib pada refresh)
	access, refresh, err := helper.GenerateToken(user.ID, roleName, []string{})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "token generation failed")
	}

	return c.JSON(fiber.Map{
		"token":   access,
		"refresh": refresh,
	})
}

//
// ------------------------------------------------------
// LOGOUT (dummy sesuai SRS)
// ------------------------------------------------------
func (s *AuthService) Logout(c *fiber.Ctx) error {

	tokenAny := c.Locals("token")
	if tokenAny == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "token not found")
	}

	token := tokenAny.(string)

	claims, err := helper.ParseToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	// blacklist sampai expired
	s.blacklist.Add(token, claims.ExpiresAt.Time)

	return c.JSON(fiber.Map{
		"message": "logout success",
	})
}

