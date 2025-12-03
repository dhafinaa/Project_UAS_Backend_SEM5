package model

import "time"

type User struct {
	ID            string    `json:"id" bson:"id"`
	Username      string    `json:"username" bson:"username"`
	Email         string    `json:"email" bson:"email"`
	Password_hash string    `json:"password_hash" bson:"password_hash"`
	Full_name     string    `json:"full_name" bson:"full_name"`
	Role_id       string    `json:"role_id" bson:"role_id"`
	Is_active     bool      `json:"is_active" bson:"is_active"`
	Created_at    time.Time `json:"created_at" bson:"created_at"`
	Updated_at    time.Time `json:"updated_at" bson:"updated_at"`
}
