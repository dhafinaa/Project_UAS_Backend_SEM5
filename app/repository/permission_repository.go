package repository

import (
	"database/sql"
)

type PermissionRepository struct {
	DB *sql.DB
}

func (r *PermissionRepository) GetByRole(roleID string) ([]string, error) {
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

	var permissions []string

	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Jika role punya 0 permission, tetap return [] (tidak error)
	return permissions, nil
}
