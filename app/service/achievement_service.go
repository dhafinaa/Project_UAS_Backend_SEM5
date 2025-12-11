package service

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

	student, _ := s.StudentRepo.FindByUserID(userID)

	ach, err := s.AchRepo.FindByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}

	if ach.Student_id != student.ID {
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

	mongoID := primitive.NewObjectID().Hex()

	achievement := model.Achievement{
		ID:               mongoID,
		Student_id:       student.ID,
		Achievement_type: req.Achievement_type,
		Title:            req.Title,
		Description:      req.Description,
		Details:          req.Details,
		Tags:             req.Tags,
		Points:           req.Points,
		Attachments:      []model.Attachment{},
		Created_at:       time.Now(),
		Updated_at:       time.Now(),
	}

	_, err = s.AchRepo.Create(c.Context(), achievement)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed saving to MongoDB")
	}

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
	if ach.Student_id != student.ID {
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
func (s *AchievementService) DeleteAchievement(c *fiber.Ctx) error {
	achID := c.Params("id")
	err := s.AchRepo.DeleteByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed deleting")
	}

	return c.JSON(fiber.Map{"message": "achievement deleted"})
}

//
// SUBMIT
//
func (s *AchievementService) SubmitAchievement(c *fiber.Ctx) error {
	achID := c.Params("id")
	userID := c.Locals("userID").(string)

	student, _ := s.StudentRepo.FindByUserID(userID)

	ach, err := s.AchRepo.FindByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}
	if ach.Student_id != student.ID {
		return fiber.NewError(fiber.StatusForbidden, "not your achievement")
	}

	err = s.AchRepo.CreateReference(c.Context(), student.ID, achID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "submit failed")
	}

	return c.JSON(fiber.Map{
		"message": "achievement submitted",
		"id":      achID,
	})
}

//
// UPLOAD ATTACHMENT
//
func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	achID := c.Params("id")
	userID := c.Locals("userID").(string)

	student, _ := s.StudentRepo.FindByUserID(userID)

	ach, err := s.AchRepo.FindByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}
	if ach.Student_id != student.ID {
		return fiber.NewError(fiber.StatusForbidden, "not your achievement")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	savePath := "./uploads/" + file.Filename
	err = c.SaveFile(file, savePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed saving file")
	}

	attachment := model.Attachment{
		File_name:  file.Filename,
		File_url:   "/uploads/" + file.Filename,
		File_type:  file.Header.Get("Content-Type"),
		Uploaded_at: time.Now(),
	}

	err = s.AchRepo.AddAttachment(c.Context(), achID, attachment)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed adding attachment")
	}

	return c.JSON(fiber.Map{
		"message":  "attachment uploaded",
		"filename": file.Filename,
	})
}
