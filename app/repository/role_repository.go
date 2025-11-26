package repository

import (
	"database/sql"
	"PROJECT_UAS/app/model"
)

type RoleRepository struct {
	DB *sql.DB
}

func (r *RoleRepository) FindByID(id string) (*model.Role, error) {
	row := r.DB.QueryRow(`
		SELECT id, name, description, created_at FROM roles WHERE id=$1
	`, id)

	var role model.Role
	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.Created_at)

	if err != nil {
		return nil, err
	}

	return &role, nil
}
