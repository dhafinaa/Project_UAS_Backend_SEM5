package repository

import (
	"database/sql"
	"PROJECT_UAS/app/model"
)

type AchievementReferenceRepository struct {
	DB *sql.DB
}

func (r *AchievementReferenceRepository) Create(ref model.Achievement_reference) error {
	_, err := r.DB.Exec(`
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`,
		ref.ID, ref.Student_id, ref.Mongo_achievement_id, ref.Status)

	return err
}

func (r *AchievementReferenceRepository) UpdateStatus(refID, status, note string) error {
	_, err := r.DB.Exec(`
		UPDATE achievement_references
		SET status=$1, rejection_note=$2, updated_at=NOW()
		WHERE id=$3
	`,
		status, note, refID)

	return err
}
