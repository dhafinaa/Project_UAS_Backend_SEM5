package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"PROJECT_UAS/app/repository"
)

type ReportService struct {
	StudentRepo     *repository.StudentRepository
	AchievementRepo *repository.AchievementRepository
	LecturerRepo    *repository.LecturerRepository
}

func NewReportService(
	sRepo *repository.StudentRepository,
	aRepo *repository.AchievementRepository,
	lRepo *repository.LecturerRepository,
) *ReportService {
	return &ReportService{
		StudentRepo:     sRepo,
		AchievementRepo: aRepo,
		LecturerRepo:    lRepo,
	}
}

//
// GET /reports/statistics
// Admin & Dosen Wali (pakai data yang ADA)
//
func (s *ReportService) GetStatistics(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats, err := s.AchievementRepo.GetStatistics(ctx)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(stats)
}

//
// GET /reports/student/:id
// Admin → bebas
// Dosen wali → mahasiswa bimbingannya
//
func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)
	studentID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ================= ADMIN =================
	if role == "Admin" {
		data, err := s.AchievementRepo.GetStudentAchievementsReport(ctx, studentID)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}
		return c.JSON(data)
	}

	// ================= DOSEN WALI =================
	lecturerID, err := s.LecturerRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(403, "lecturer not found")
	}

	// cek apakah student ini bimbingannya
	students, err := s.StudentRepo.FindByAdvisor(lecturerID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	isAdvisor := false
	for _, st := range students {
		if st.ID == studentID {
			isAdvisor = true
			break
		}
	}

	if !isAdvisor {
		return fiber.NewError(403, "not your student")
	}

	data, err := s.AchievementRepo.GetStudentAchievementsReport(ctx, studentID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(data)
}