package service

import "github.com/gofiber/fiber/v2"

type LecturerService struct{}

func NewLecturerService() *LecturerService {
	return &LecturerService{}
}

// GET /lecturer/achievements
func (s *LecturerService) GetStudentAchievements(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "List student achievements (dummy)"})
}

// PUT /lecturer/achievements/:id/verify
func (s *LecturerService) VerifyAchievement(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Verify achievement (dummy)"})
}

// PUT /lecturer/achievements/:id/reject
func (s *LecturerService) RejectAchievement(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Reject achievement (dummy)"})
}