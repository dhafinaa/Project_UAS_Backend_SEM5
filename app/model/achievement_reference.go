package model

import "time"

type Achievement_reference struct {
	ID                 string    `json:"id"`
	Student_id         string    `json:"student_id"`
	Mongo_achievement_id string  `json:"mongo_achievement_id"`
	Status             string    `json:"status"`
	Submitted_at       time.Time `json:"submitted_at"`
	Verified_at        time.Time `json:"verified_at"`
	Verified_by        string    `json:"verified_by"`
	Rejection_note     string    `json:"rejection_note"`
	Created_at         time.Time `json:"created_at"`
	Updated_at         time.Time `json:"updated_at"`
}
