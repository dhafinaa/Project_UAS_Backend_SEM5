package repository

import (
	"database/sql"
	"errors"

	"PROJECT_UAS/app/model"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(pg *sql.DB) *AuthRepository {
	return &AuthRepository{DB: pg}
}

// =======================================================
// FIND USER BY USERNAME OR EMAIL
// =======================================================
func (r *AuthRepository) FindByLogin(login string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active
		FROM users
		WHERE username = $1 OR email = $1
	`

	row := r.DB.QueryRow(query, login)

	var u model.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.Password_hash,
		&u.Full_name, &u.Role_id, &u.Is_active,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}

// =======================================================
// GET PERMISSIONS BY ROLE ID
// =======================================================
func (r *AuthRepository) GetPermissionsByRole(roleID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := r.DB.Query(query, roleID)
	if err != nil {
		return nil, errors.New("query failed: " + err.Error())
	}
	defer rows.Close()

	var perms []string

	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, errors.New("scan failed: " + err.Error())
		}
		perms = append(perms, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return perms, nil
}
