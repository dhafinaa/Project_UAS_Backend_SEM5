package service

import (
	"github.com/gofiber/fiber/v2"

	"PROJECT_UAS/app/repository"
)

type StudentService struct {
	studentRepo  *repository.StudentRepository
	lecturerRepo *repository.LecturerRepository
}

func NewStudentService(
	studentRepo *repository.StudentRepository,lecturerRepo *repository.LecturerRepository,) *StudentService {
	return &StudentService{
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
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
