package service

import (
	"errors"
	"PROJECT_UAS/app/model"
	"PROJECT_UAS/app/repository"
)

type StudentService struct {
	studentRepo *repository.StudentRepository
}

func NewStudentService(studentRepo *repository.StudentRepository) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
	}
}

// =============================
// GET STUDENT BY ID
// =============================
func (s *StudentService) GetByID(id string) (*model.Student, error) {
	if id == "" {
		return nil, errors.New("id tidak boleh kosong")
	}

	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return student, nil
}

// =============================
// GET STUDENT BY USER ID
// =============================
func (s *StudentService) GetByUserID(userID string) (*model.Student, error) {
	if userID == "" {
		return nil, errors.New("user_id tidak boleh kosong")
	}

	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return student, nil
}

// =============================
// GET ALL STUDENTS BY ADVISOR / LECTURER
// =============================
func (s *StudentService) GetStudentsByAdvisor(lecturerID string) ([]model.Student, error) {
	if lecturerID == "" {
		return nil, errors.New("advisor_id tidak boleh kosong")
	}

	students, err := s.studentRepo.FindByAdvisor(lecturerID)
	if err != nil {
		return nil, err
	}

	return students, nil
}
