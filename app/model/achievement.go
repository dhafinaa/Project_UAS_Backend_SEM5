package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	StudentID       string                 `bson:"student_id" json:"student_id"`
	AchievementType string                 `bson:"achievement_type" json:"achievement_type"`
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`

	Details     map[string]interface{} `bson:"details" json:"details"`
	Attachments []Attachment           `bson:"attachments" json:"attachments"`

	Tags   []string `bson:"tags" json:"tags"`
	Points int      `bson:"points" json:"points"`

	IsDeleted        bool      `bson:"is_deleted" json:"is_deleted"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Attachment struct {
	FileName   string    `bson:"file_name" json:"file_name"`
	FileURL    string    `bson:"file_url" json:"file_url"`
	FileType   string    `bson:"file_type" json:"file_type"`
	UploadedAt time.Time `bson:"uploaded_at" json:"uploaded_at"`
}
