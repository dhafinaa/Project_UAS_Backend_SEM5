package service

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"PROJECT_UAS/app/model"
	"PROJECT_UAS/app/repository"
)

type StudentService struct {
	AchRepo     *repository.AchievementRepository
	StudentRepo *repository.StudentRepository
}

func NewStudentService(achRepo *repository.AchievementRepository, studentRepo *repository.StudentRepository) *StudentService {
	return &StudentService{
		AchRepo:     achRepo,
		StudentRepo: studentRepo,
	}
}


// ----------------------------------------------------------------------
// GET ACHIEVEMENTS (Mahasiswa)
// ----------------------------------------------------------------------
func (s *StudentService) GetAchievements(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student profile not found")
	}

	list, err := s.AchRepo.FindByStudentID(c.Context(), student.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed loading achievements")
	}

	return c.JSON(list)
}

// ----------------------------------------------------------------------
// CREATE ACHIEVEMENT
// ----------------------------------------------------------------------
func (s *StudentService) CreateAchievement(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student profile not found")
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

	newID := primitive.NewObjectID().Hex()

	achievement := model.Achievement{
		ID:               newID,
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
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save achievement")
	}

	return c.JSON(fiber.Map{
		"message":       "achievement created",
		"achievement_id": newID,
	})
}

// ----------------------------------------------------------------------
// SUBMIT ACHIEVEMENT — membuat record reference di SQL
// ----------------------------------------------------------------------
func (s *StudentService) SubmitAchievement(c *fiber.Ctx) error {

	achID := c.Params("id")
	userID := c.Locals("userID").(string)

	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student profile not found")
	}

	ref := model.Achievement_reference{
		ID:                 primitive.NewObjectID().Hex(),
		Student_id:         student.ID,
		Mongo_achievement_id: achID,
		Status:             "submitted",
		Submitted_at:       time.Now(),
		Created_at:         time.Now(),
		Updated_at:         time.Now(),
	}

	err = s.AchRepo.CreateReference(c.Context(), ref)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed submitting achievement")
	}

	return c.JSON(fiber.Map{
		"message": "achievement submitted",
		"ref_id":  ref.ID,
	})
}

// ----------------------------------------------------------------------
// DELETE ACHIEVEMENT
// ----------------------------------------------------------------------
func (s *StudentService) DeleteAchievement(c *fiber.Ctx) error {

	achID := c.Params("id")

	err := s.AchRepo.DeleteByID(c.Context(), achID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed deleting achievement")
	}

	return c.JSON(fiber.Map{
		"message": "achievement deleted",
		"id":      achID,
	})
}
