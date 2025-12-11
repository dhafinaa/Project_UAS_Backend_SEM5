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

	// Ambil student profile berdasarkan user_id
	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student profile not found")
	}

	// Request body
	var req struct {
		Achievement_type string                 `json:"achievement_type"`
		Title            string                 `json:"title"`
		Description      string                 `json:"description"`
		Details          map[string]interface{} `json:"details"`
		Tags             []string               `json:"tags"`
		Points           int                    `json:"points"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input format")
	}

	// ID MongoDB manual (Hex string)
	mongoID := primitive.NewObjectID().Hex()

	// Siapkan data prestasi untuk MongoDB
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

	// ------------------------------
	// 1) INSERT ke MongoDB
	// ------------------------------
	_, err = s.AchRepo.Create(c.Context(), achievement)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save achievement to MongoDB")
	}

	// ------------------------------
	// 2) INSERT reference ke PostgreSQL
	// Status awal = "draft"
	// submitted_at = NULL
	// ------------------------------
	err = s.AchRepo.CreateReference(c.Context(), student.ID, mongoID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save reference to PostgreSQL")
	}

	// Response SRS FR-003
	return c.JSON(fiber.Map{
		"message":           "achievement created",
		"achievement_id":    mongoID,
		"reference_status":  "draft",
	})
}



// ----------------------------------------------------------------------
// SUBMIT ACHIEVEMENT — membuat record reference di SQL
// ----------------------------------------------------------------------
// SUBMIT ACHIEVEMENT
func (s *StudentService) SubmitAchievement(c *fiber.Ctx) error {

	achID := c.Params("id")
	userID := c.Locals("userID").(string)

	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "student profile not found")
	}

	ctx := c.Context()

	// Verify achievement exists
	ach, err := s.AchRepo.FindByID(ctx, achID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "achievement not found")
	}
	if ach.Student_id != student.ID {
		return fiber.NewError(fiber.StatusForbidden, "achievement does not belong to student")
	}

	// Insert reference into PostgreSQL with status = 'submitted'
	err = s.AchRepo.CreateReference(ctx, student.ID, achID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed submitting achievement")
	}

	return c.JSON(fiber.Map{
		"message": "achievement submitted",
		"achievement_id": achID,
		"status": "submitted",
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
