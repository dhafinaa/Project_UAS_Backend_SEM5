package model

import "time"

type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password_hash string    `json:"password_hash"`
	Full_name     string    `json:"full_name"`
	Role_id       string    `json:"role_id"`
	Is_active     bool      `json:"is_active"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}
