package service

import (
    "PROJECT_UAS/app/repository"

	 "github.com/gofiber/fiber/v2"
)

type ReportService struct {
    AchRepo *repository.AchievementRepository
}

func NewReportService(achRepo *repository.AchievementRepository) *ReportService {
    return &ReportService{AchRepo: achRepo}
}

func (s *ReportService) GetStatistics(c *fiber.Ctx) error {

    stats, err := s.AchRepo.GetStatistics(c.Context())
    if err != nil {
        return fiber.NewError(500, err.Error())
    }

    return c.JSON(stats)
}


func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {

    studentID := c.Params("id")

    achievements, err := s.AchRepo.
        GetStudentAchievementsReport(c.Context(), studentID)

    if err != nil {
        return fiber.NewError(500, err.Error())
    }

    return c.JSON(achievements)
}
