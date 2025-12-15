package service

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"PROJECT_UAS/app/model"
	"PROJECT_UAS/app/repository"
)

type AchievementService struct {
	AchRepo     *repository.AchievementRepository
	StudentRepo *repository.StudentRepository
}

func NewAchievementService(achRepo *repository.AchievementRepository, studentRepo *repository.StudentRepository) *AchievementService {
	return &AchievementService{
		AchRepo:     achRepo,
		StudentRepo: studentRepo,
	}
}

//
// GET ALL ACHIEVEMENTS (Mahasiswa)
//
func (s *AchievementService) GetAchievements(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student not found")
	}

	list, err := s.AchRepo.FindByStudentID(c.Context(), student.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed loading achievements")
	}

	return c.JSON(list)
}

//
// GET DETAIL
//
func (s *AchievementService) GetAchievementDetail(c *fiber.Ctx) error {
	achID := c.Params("id")
	userID := c.Locals("userID").(string)

	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student not found")
	}

	ach, err := s.AchRepo.FindByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}

	if ach.StudentID != student.ID {
		return fiber.NewError(fiber.StatusForbidden, "not your achievement")
	}

	return c.JSON(ach)
}


//
// CREATE
//
func (s *AchievementService) CreateAchievement(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student not found")
	}

	var req struct {
		Achievement_type string                 `json:"achievement_type"`
		Title            string                 `json:"title"`
		Description      string                 `json:"description"`
		Details          map[string]interface{} `json:"details"`
		Tags             []string               `json:"tags"`
		Points           int                    `json:"points"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	achievement := model.Achievement{
		StudentID:       student.ID,
		AchievementType: req.Achievement_type,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
		Attachments:     []model.Attachment{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 1️⃣ INSERT ke Mongo
	mongoID, err := s.AchRepo.Create(c.Context(), achievement)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed saving to MongoDB")
	}

	// 2️⃣ INSERT reference ke PostgreSQL
	err = s.AchRepo.CreateReference(c.Context(), student.ID, mongoID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed saving reference")
	}

	return c.JSON(fiber.Map{
		"message":         "achievement created",
		"achievement_id":  mongoID,
		"reference_status": "draft",
	})
}


//
// UPDATE
//
func (s *AchievementService) UpdateAchievement(c *fiber.Ctx) error {
	achID := c.Params("id")
	userID := c.Locals("userID").(string)

	student, _ := s.StudentRepo.FindByUserID(userID)

	ach, err := s.AchRepo.FindByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}
	if ach.StudentID != student.ID {
		return fiber.NewError(fiber.StatusForbidden, "not your achievement")
	}

	var req struct {
		Title       string                 `json:"title"`
		Description string                 `json:"description"`
		Details     map[string]interface{} `json:"details"`
		Tags        []string               `json:"tags"`
		Points      int                    `json:"points"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	update := bson.M{
		"title":       req.Title,
		"description": req.Description,
		"details":     req.Details,
		"tags":        req.Tags,
		"points":      req.Points,
		"updated_at":  time.Now(),
	}

	err = s.AchRepo.UpdateAchievement(c.Context(), achID, update)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed updating achievement")
	}

	return c.JSON(fiber.Map{"message": "achievement updated"})
}

//
// DELETE
//
// func (s *AchievementService) DeleteAchievement(c *fiber.Ctx) error {
// 	achID := c.Params("id")
// 	err := s.AchRepo.DeleteByID(c.Context(), achID)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "failed deleting")
// 	}

// 	return c.JSON(fiber.Map{"message": "achievement deleted"})
// }

//
// SUBMIT
//
func (s *AchievementService) SubmitAchievement(c *fiber.Ctx) error {
	achID := c.Params("id")
	userID := c.Locals("userID").(string)

	// 1. Ambil student
	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student not found")
	}

	// 2. Pastikan achievement ada di Mongo
	ach, err := s.AchRepo.FindByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}

	// 3. Pastikan milik mahasiswa
	if ach.StudentID != student.ID {
		return fiber.NewError(fiber.StatusForbidden, "not your achievement")
	}

	// 4. Update status di PostgreSQL
	err = s.AchRepo.SubmitAchievement(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"message":        "achievement submitted",
		"achievement_id": achID,
		"status":         "submitted",
	})
}


//
// UPLOAD ATTACHMENT
//
func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	achID := c.Params("id")
	userID := c.Locals("userID").(string)

	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student not found")
	}

	ach, err := s.AchRepo.FindByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}

	if ach.StudentID != student.ID {
		return fiber.NewError(fiber.StatusForbidden, "not your achievement")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	savePath := "./uploads/" + file.Filename
	if err := c.SaveFile(file, savePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed saving file")
	}

	attachment := model.Attachment{
		FileName:   file.Filename,
		FileURL:    "/uploads/" + file.Filename,
		FileType:   file.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	if err := s.AchRepo.AddAttachment(c.Context(), achID, attachment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed adding attachment")
	}

	return c.JSON(fiber.Map{
		"message":  "attachment uploaded",
		"filename": file.Filename,
	})
}

//soft delete
func (s *AchievementService) DeleteAchievement(c *fiber.Ctx) error {
	achID := c.Params("id") // mongo_achievement_id
	userID := c.Locals("userID").(string)

	// 1. Ambil student dari token
	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student not found")
	}

	// 2. Validasi achievement ada di Mongo
	ach, err := s.AchRepo.FindByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}

	// 3. Validasi kepemilikan
	if ach.StudentID != student.ID {
		return fiber.NewError(fiber.StatusForbidden, "not your achievement")
	}

	// 4. Update PostgreSQL → hanya draft boleh delete
	err = s.AchRepo.DeleteDraftAchievement(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"message":        "achievement deleted",
		"achievement_id": achID,
		"status":         "deleted",
	})
}
