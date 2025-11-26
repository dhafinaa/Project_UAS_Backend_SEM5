package model

import "time"

type Lecturer struct {
	ID          string    `json:"id"`
	User_id     string    `json:"user_id"`
	Lecturer_id string    `json:"lecturer_id"`
	Department  string    `json:"department"`
	Created_at  time.Time `json:"created_at"`
}
