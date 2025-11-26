package repository

import (
	"database/sql"
	"PROJECT_UAS/app/model"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE id=$1
	`, id)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password_hash,
		&u.Full_name, &u.Role_id, &u.Is_active, &u.Created_at, &u.Updated_at)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE username=$1 OR email=$1 LIMIT 1
	`, login)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password_hash,
		&u.Full_name, &u.Role_id, &u.Is_active, &u.Created_at, &u.Updated_at)

	if err != nil {
		return nil, err
	}
	return &u, nil
}
