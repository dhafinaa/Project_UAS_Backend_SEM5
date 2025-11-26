package repository

import (
	"database/sql"
)

type PermissionRepository struct {
	DB *sql.DB
}

func (r *PermissionRepository) GetByRole(roleID string) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT p.name 
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`, roleID)

	if err != nil {
		return nil, err
	}

	var permissions []string
	for rows.Next() {
		var perm string
		rows.Scan(&perm)
		permissions = append(permissions, perm)
	}

	return permissions, nil
}
