package repository

import (
	"database/sql"
	"PROJECT_UAS/app/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GET ALL USERS
func (r *UserRepository) FindAll() ([]model.User, error) {
	rows, err := r.DB.Query(`
		SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Full_name,
			&u.Role_id,
			&u.Is_active,
			&u.Created_at,
			&u.Updated_at,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GET USER BY ID
func (r *UserRepository) FindByID(id string) (*model.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE id=$1
	`, id)

	var u model.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Full_name,
		&u.Role_id,
		&u.Is_active,
		&u.Created_at,
		&u.Updated_at,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CREATE USER
func (r *UserRepository) Create(u model.User) error {
	_, err := r.DB.Exec(`
		INSERT INTO users
		(username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,true,NOW(),NOW())
	`,
		u.Username,
		u.Email,
		u.Password_hash,
		u.Full_name,
		u.Role_id,
	)
	return err
}

// UPDATE USER (NO PASSWORD)
func (r *UserRepository) Update(id string, u model.User) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET username=$1,
		    email=$2,
		    full_name=$3,
		    updated_at=NOW()
		WHERE id=$4
	`,
		u.Username,
		u.Email,
		u.Full_name,
		id,
	)
	return err
}

// SOFT DELETE USER
func (r *UserRepository) SoftDelete(id string) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET is_active=false, updated_at=NOW()
		WHERE id=$1
	`, id)
	return err
}

// UPDATE ROLE USER
func (r *UserRepository) UpdateRole(userID, roleID string) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET role_id=$1, updated_at=NOW()
		WHERE id=$2
	`, roleID, userID)
	return err
}
