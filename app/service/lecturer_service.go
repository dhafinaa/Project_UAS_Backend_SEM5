package service

import (

	"PROJECT_UAS/app/model"
	"PROJECT_UAS/app/repository"

	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	StudentRepo     *repository.StudentRepository
	AchievementRepo *repository.AchievementRepository
	LecturerRepo  *repository.LecturerRepository
}

func NewLecturerService(sRepo *repository.StudentRepository, aRepo *repository.AchievementRepository, lRepo *repository.LecturerRepository) *LecturerService {
	return &LecturerService{
		StudentRepo:     sRepo,
		AchievementRepo: aRepo,
		LecturerRepo: lRepo,
	}
}

// GET /lecturer/advisees
func (s *LecturerService) GetStudentAchievements(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// 1. Ambil lecturer.id dari users.id
	lecturerID, err := s.LecturerRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(404, "lecturer not found")
	}

	// 2. Ambil mahasiswa bimbingan
	students, err := s.StudentRepo.FindByAdvisor(lecturerID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	if len(students) == 0 {
		return c.JSON([]model.Achievement{})
	}

	// 3. Ambil ID mahasiswa
	var studentIDs []string
	for _, st := range students {
		studentIDs = append(studentIDs, st.ID)
	}

	// 4. Ambil achievement SUBMITTED
	achievements, err := s.AchievementRepo.
		ListSubmittedByStudents(c.Context(), studentIDs)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(achievements)
}



// PUT /lecturer/achievements/:id/verify
func (s *LecturerService) VerifyAchievement(c *fiber.Ctx) error {
	achID := c.Params("id")
	userID := c.Locals("userID").(string)

	// 1. Ambil lecturer_id dari user_id
	lecturerID, err := s.LecturerRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "lecturer not found")
	}

	// 2. Verify achievement (PostgreSQL only)
	err = s.AchievementRepo.VerifyAchievement(
		c.Context(),
		achID,
		lecturerID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"message":        "achievement verified",
		"achievement_id": achID,
		"status":         "verified",
	})
}


// PUT /lecturer/achievements/:id/reject
func (s *LecturerService) RejectAchievement(c *fiber.Ctx) error {
    achID := c.Params("id")
    userID := c.Locals("userID").(string)

    // 1. Ambil lecturer_id dari user_id
    lecturerID, err := s.LecturerRepo.FindByUserID(userID)
    if err != nil {
        return fiber.NewError(403, "lecturer not found")
    }

    // 2. Parse request body (rejection note)
    var req struct {
        Note string `json:"rejection_note"`
    }

    if err := c.BodyParser(&req); err != nil || req.Note == "" {
        return fiber.NewError(400, "rejection_note is required")
    }

    // 3. Update status → rejected
    err = s.AchievementRepo.RejectAchievement(
        c.Context(),
        achID,
        lecturerID,
        req.Note,
    )
    if err != nil {
        return fiber.NewError(400, err.Error())
    }

    // 4. (Optional) Notification hook
    // notificationService.NotifyStudentReject(...)

    return c.JSON(fiber.Map{
        "message":        "achievement rejected",
        "achievement_id": achID,
        "status":         "rejected",
    })
}


