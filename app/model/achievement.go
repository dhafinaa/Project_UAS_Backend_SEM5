package model

import "time"

type Achievement struct {
	ID                string                 `json:"_id"`
	Student_id        string                 `json:"student_id"`
	Achievement_type  string                 `json:"achievement_type"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`

	Details           map[string]interface{} `json:"details"`

	Attachments       []Attachment           `json:"attachments"`

	Tags              []string               `json:"tags"`
	Points            int                    `json:"points"`

	Created_at        time.Time              `json:"created_at"`
	Updated_at        time.Time              `json:"updated_at"`
}

type Attachment struct {
	File_name  string    `json:"file_name"`
	File_url   string    `json:"file_url"`
	File_type  string    `json:"file_type"`
	Uploaded_at time.Time `json:"uploaded_at"`
}
