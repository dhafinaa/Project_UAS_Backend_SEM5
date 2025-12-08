package service

import (
	"context"
	"time"

	"PROJECT_UAS/app/model"
	"PROJECT_UAS/app/repository"

	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	StudentRepo     *repository.StudentRepository
	AchievementRepo *repository.AchievementRepository
}

func NewLecturerService(sRepo *repository.StudentRepository, aRepo *repository.AchievementRepository) *LecturerService {
	return &LecturerService{
		StudentRepo:     sRepo,
		AchievementRepo: aRepo,
	}
}

// GET /lecturer/advisees
func (s *LecturerService) GetStudentAchievements(c *fiber.Ctx) error {

	lecturerID := c.Locals("userID")
	if lecturerID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "no user")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// load mahasiswa bimbingan
	students, err := s.StudentRepo.FindByAdvisor(lecturerID.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := make(map[string][]model.Achievement)

	for _, st := range students {
		achievements, _ := s.AchievementRepo.ListByStudent(ctx, st.ID)
		response[st.ID] = achievements
	}

	return c.JSON(response)
}

// PUT /lecturer/achievements/:id/verify
func (s *LecturerService) VerifyAchievement(c *fiber.Ctx) error {

	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.AchievementRepo.UpdateStatusByID(ctx, id, "verified")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{"message": "Achievement verified"})
}

// PUT /lecturer/achievements/:id/reject
func (s *LecturerService) RejectAchievement(c *fiber.Ctx) error {

	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.AchievementRepo.UpdateStatusByID(ctx, id, "rejected")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{"message": "Achievement rejected"})
}
