package service

import (
	"github.com/gofiber/fiber/v2"

	"PROJECT_UAS/app/model"
	"PROJECT_UAS/app/repository"
	"PROJECT_UAS/helper"
)

type AuthService struct {
	Repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	// Find user
	user, err := s.Repo.FindByLogin(req.Login)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid username/email or password")
	}

	// Check active status
	if !user.Is_active {
		return fiber.NewError(fiber.StatusUnauthorized, "user is inactive")
	}

	// Validate password
	if !helper.CheckPasswordHash(req.Password, user.Password_hash) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid username/email or password")
	}

	// Load permissions
	perms, err := s.Repo.GetPermissionsByRole(user.Role_id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed load permissions")
	}

	// Generate JWT
	token, refresh, err := helper.GenerateToken(user.ID, user.Role_id, perms)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed generating token")
	}

	// Build response
	resp := model.LoginResponse{
		Token:       token,
		Refresh:     refresh,
		User:        *user,
		Permissions: perms,
	}

	return c.JSON(resp)
}

// Not implemented
func (s *AuthService) RefreshToken(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "refresh token not implemented"})
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "logout not implemented"})
}

func (s *AuthService) Profile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "profile not implemented"})
}
