package service

import (
	"PROJECT_UAS/app/model"
	"PROJECT_UAS/app/repository"
	"PROJECT_UAS/helper"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GET ALL USERS
func (s *UserService) GetUsers(c *fiber.Ctx) error {
	users, err := s.repo.FindAll()
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	return c.JSON(users)
}

// GET USER DETAIL
func (s *UserService) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := s.repo.FindByID(id)
	if err != nil {
		return fiber.NewError(404, "user not found")
	}
	return c.JSON(user)
}

// CREATE USER
func (s *UserService) CreateUser(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		RoleID   string `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid input")
	}

	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(500, "password hashing failed")
	}

	user := model.User{
		Username:      req.Username,
		Email:         req.Email,
		Password_hash: hash,
		Full_name:     req.FullName,
		Role_id:       req.RoleID,
	}

	if err := s.repo.Create(user); err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "user created successfully",
	})
}

// UPDATE USER
func (s *UserService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid input")
	}

	user := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Full_name: req.FullName,
	}

	if err := s.repo.Update(id, user); err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "user updated successfully",
	})
}

// DELETE USER (SOFT)
func (s *UserService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.repo.SoftDelete(id); err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "user deactivated",
	})
}

// UPDATE USER ROLE
func (s *UserService) UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		RoleID string `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid input")
	}

	if err := s.repo.UpdateRole(id, req.RoleID); err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "user role updated",
	})
}