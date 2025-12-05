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

// ------------------------------------------------------
// LOGIN — mencari user berdasarkan username/email
// ------------------------------------------------------
func (r *AuthRepository) FindByLogin(login string) (*model.User, string, error) {

	query := `
		SELECT u.id, u.username, u.email, u.password_hash, 
		       u.full_name, u.role_id, u.is_active,
		       r.name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.username = $1 OR u.email = $1
	`

	row := r.DB.QueryRow(query, login)

	var u model.User
	var roleName string

	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.Password_hash,
		&u.Full_name, &u.Role_id, &u.Is_active,
		&roleName,
	)

	if err != nil {
		return nil, "", errors.New("user not found")
	}

	return &u, roleName, nil
}

// ------------------------------------------------------
// LOAD USER + ROLE FULL (untuk /profile)
// ------------------------------------------------------
func (r *AuthRepository) GetUserRoleByID(userID string) (*model.User, string, []string, error) {

	query := `
		SELECT u.id, u.username, u.email, u.password_hash,
		       u.full_name, u.role_id, u.is_active,
		       r.name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`

	row := r.DB.QueryRow(query, userID)

	var user model.User
	var roleName string

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password_hash,
		&user.Full_name, &user.Role_id, &user.Is_active,
		&roleName,
	)

	if err != nil {
		return nil, "", nil, err
	}

	// Load permissions
	perms, err := r.GetPermissionsByRoleID(user.Role_id)
	if err != nil {
		return &user, roleName, []string{}, nil
	}

	return &user, roleName, perms, nil
}

// ------------------------------------------------------
// LOAD PERMISSIONS BY ROLE ID
// ------------------------------------------------------
func (r *AuthRepository) GetPermissionsByRoleID(roleID string) ([]string, error) {

	query := `
		SELECT p.name
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`

	rows, err := r.DB.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string

	for rows.Next() {
		var p string
		rows.Scan(&p)
		perms = append(perms, p)
	}

	return perms, nil
}

// GetPermissionsByRoleName returns permissions for a role by its name (SRS-FR002)
func (r *AuthRepository) GetPermissionsByRoleName(roleName string) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN roles r2 ON r2.id = rp.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE LOWER(r2.name) = LOWER($1)
	`

	rows, err := r.DB.Query(query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}

	return perms, nil
}
