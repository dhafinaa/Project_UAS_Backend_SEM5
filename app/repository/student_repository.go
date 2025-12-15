package repository

import (
	"database/sql"
	"PROJECT_UAS/app/model"
)

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) FindByID(id string) (*model.Student, error) {

	row := r.DB.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students 
		WHERE id=$1
	`, id)

	var s model.Student
	err := row.Scan(&s.ID, &s.User_id, &s.Student_id,
		&s.Program_study, &s.Academic_year, &s.Advisor_id, &s.Created_at)

	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) FindByUserID(userID string) (*model.Student, error) {

	row := r.DB.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students 
		WHERE user_id=$1
	`, userID)

	var s model.Student
	err := row.Scan(&s.ID, &s.User_id, &s.Student_id,
		&s.Program_study, &s.Academic_year, &s.Advisor_id, &s.Created_at)

	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) FindByAdvisor(lecturerID string) ([]model.Student, error) {

	rows, err := r.DB.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students 
		WHERE advisor_id=$1
	`, lecturerID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Student

	for rows.Next() {
		var s model.Student
		rows.Scan(&s.ID, &s.User_id, &s.Student_id,
			&s.Program_study, &s.Academic_year, &s.Advisor_id, &s.Created_at)

		list = append(list, s)
	}

	return list, nil
}


// Digunakan oleh Dosen Wali
func (r *StudentRepository) FindByAdvisorID(advisorID string) ([]string, error) {
	query := `
		SELECT id
		FROM students
		WHERE advisor_id = $1
	`

	rows, err := r.DB.Query(query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentIDs []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		studentIDs = append(studentIDs, id)
	}

	return studentIDs, nil
}


