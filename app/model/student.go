package model

import "time"

type Student struct {
	ID            string    `json:"id"`
	User_id       string    `json:"user_id"`
	Student_id    string    `json:"student_id"`
	Program_study string    `json:"program_study"`
	Academic_year string    `json:"academic_year"`
	Advisor_id    string    `json:"advisor_id"`
	Created_at    time.Time `json:"created_at"`
}
