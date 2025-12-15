package repository

import (
	"database/sql"

	"PROJECT_UAS/app/model"
)

type LecturerRepository struct {
	DB *sql.DB
}

func NewLecturerRepository(db *sql.DB) *LecturerRepository {
	return &LecturerRepository{
		DB: db,
	}
}

// Cari lecturer berdasarkan PK lecturers.id
func (r *LecturerRepository) FindByID(id string) (*model.Lecturer, error) {
	row := r.DB.QueryRow(`
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE id = $1
	`, id)

	var l model.Lecturer
	err := row.Scan(
		&l.ID,
		&l.User_id,
		&l.Lecturer_id,
		&l.Department,
		&l.Created_at,
	)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

// 🔥 INI YANG PALING PENTING UNTUK FR-006
// Mapping dari users.id (JWT) → lecturers.id
func (r *LecturerRepository) FindByUserID(userID string) (string, error) {
	query := `
		SELECT id
		FROM lecturers
		WHERE user_id = $1
	`

	var lecturerID string
	err := r.DB.QueryRow(query, userID).Scan(&lecturerID)
	if err != nil {
		return "", err
	}

	return lecturerID, nil
}
