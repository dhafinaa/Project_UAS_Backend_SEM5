package repository

import (
	"database/sql"
	"PROJECT_UAS/app/model"
)

type LecturerRepository struct {
	DB *sql.DB
}

func (r *LecturerRepository) FindByID(id string) (*model.Lecturer, error) {
	row := r.DB.QueryRow(`
		SELECT id, user_id, lecturer_id, department, created_at 
		FROM lecturers WHERE id=$1
	`, id)

	var l model.Lecturer
	err := row.Scan(&l.ID, &l.User_id, &l.Lecturer_id,
		&l.Department, &l.Created_at)

	if err != nil {
		return nil, err
	}

	return &l, nil
}
