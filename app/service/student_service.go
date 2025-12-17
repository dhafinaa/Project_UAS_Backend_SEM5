package service

import (
	"github.com/gofiber/fiber/v2"

	"PROJECT_UAS/app/repository"
)

type StudentService struct {
	studentRepo  *repository.StudentRepository
	lecturerRepo *repository.LecturerRepository
	achievementRepo *repository.AchievementRepository

}

func NewStudentService(
	studentRepo *repository.StudentRepository, lecturerRepo *repository.LecturerRepository, achievementRepo *repository.AchievementRepository) *StudentService {
	return &StudentService{
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
		achievementRepo: achievementRepo,
	}
}


// =============================
// GET ALL STUDENTS (ADMIN)
// =============================
func (s *StudentService) GetAllStudents(c *fiber.Ctx) error {

	students, err := s.studentRepo.FindAll()
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to load students",
		)
	}

	return c.JSON(students)
}

// =============================
// GET STUDENT BY ID
// =============================
func (s *StudentService) GetStudentByID(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"id is required",
		)
	}

	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			"student not found",
		)
	}

	return c.JSON(student)
}

// =============================
// GET STUDENTS BY ADVISOR (DOSEN WALI)
// =============================
func (s *StudentService) GetStudentsByAdvisor(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	// 🔥 STEP 1: users.id → lecturers.id
	lecturerID, err := s.lecturerRepo.FindByUserID(userID)
	if err != nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			"lecturer not found",
		)
	}

	// 🔥 STEP 2: pakai lecturers.id
	students, err := s.studentRepo.FindByAdvisor(lecturerID)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed loading students",
		)
	}

	return c.JSON(students)
}

func (s *StudentService) GetStudentAchievements(c *fiber.Ctx) error {

    studentIDParam := c.Params("id")
    if studentIDParam == "" {
        return fiber.NewError(fiber.StatusBadRequest, "student id required")
    }

    userID := c.Locals("userID").(string)

    student, err := s.studentRepo.FindByUserID(userID)
    if err != nil {
        return fiber.NewError(fiber.StatusNotFound, "student not found")
    }

    // 🔒 mahasiswa hanya boleh akses data sendiri
    if student.ID != studentIDParam {
        return fiber.NewError(
            fiber.StatusForbidden,
            "cannot access other student's achievements",
        )
    }

    achievements, err := s.achievementRepo.
        FindByStudentID(c.Context(), student.ID)
    if err != nil {
        return fiber.NewError(
            fiber.StatusInternalServerError,
            "failed loading achievements",
        )
    }

    return c.JSON(achievements)
}

func (s *StudentService) UpdateStudentAdvisor(c *fiber.Ctx) error {

    studentID := c.Params("id")
    if studentID == "" {
        return fiber.NewError(
            fiber.StatusBadRequest,
            "student id is required",
        )
    }

    var req struct {
        AdvisorID string `json:"advisor_id"`
    }

    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(
            fiber.StatusBadRequest,
            "invalid request body",
        )
    }

    if req.AdvisorID == "" {
        return fiber.NewError(
            fiber.StatusBadRequest,
            "advisor_id is required",
        )
    }

    // pastikan student ada
    _, err := s.studentRepo.FindByID(studentID)
    if err != nil {
        return fiber.NewError(
            fiber.StatusNotFound,
            "student not found",
        )
    }

    // update advisor
    err = s.studentRepo.UpdateAdvisor(studentID, req.AdvisorID)
    if err != nil {
        return fiber.NewError(
            fiber.StatusInternalServerError,
            "failed updating advisor",
        )
    }

    return c.JSON(fiber.Map{
        "message":    "advisor updated successfully",
        "student_id": studentID,
        "advisor_id": req.AdvisorID,
    })
}
