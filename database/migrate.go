package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"PROJECT_UAS/helper"
)

// RunMigration dipanggil dari main.go
func RunMigration(pg *sql.DB) {

	// ===============================
	// 1. CREATE TABLES
	// ===============================

	_, err := pg.Exec(`
	CREATE TABLE IF NOT EXISTS roles (
		id VARCHAR(50) PRIMARY KEY,
		name VARCHAR(50) UNIQUE NOT NULL,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(50) PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		full_name VARCHAR(100),
		role_id VARCHAR(50) REFERENCES roles(id),
		is_active BOOLEAN DEFAULT TRUE
	);

	CREATE TABLE IF NOT EXISTS lecturers (
		id VARCHAR(50) PRIMARY KEY,
		user_id VARCHAR(50) UNIQUE REFERENCES users(id),
		lecturer_id VARCHAR(50) UNIQUE NOT NULL,
		department VARCHAR(100)
	);

	CREATE TABLE IF NOT EXISTS students (
		id VARCHAR(50) PRIMARY KEY,
		user_id VARCHAR(50) UNIQUE REFERENCES users(id),
		student_id VARCHAR(50) UNIQUE NOT NULL,
		program_study VARCHAR(100),
		academic_year VARCHAR(10),
		advisor_id VARCHAR(50) REFERENCES lecturers(id)
	);
	`)
	if err != nil {
		log.Fatal("❌ Error creating tables:", err)
	}

	log.Println("✔ Tables created")

	// ===============================
	// 2. INSERT ROLES
	// ===============================

	roleAdmin := "ROLE-ADMIN"
	roleLect := "ROLE-LECTURER"
	roleStud := "ROLE-STUDENT"

	_, err = pg.Exec(`
	INSERT INTO roles (id, name, description) VALUES
		($1, 'Admin', 'Administrator'),
		($2, 'Dosen Wali', 'Academic Advisor'),
		($3, 'Mahasiswa', 'Student User')
	ON CONFLICT (name) DO NOTHING;
	`, roleAdmin, roleLect, roleStud)

	if err != nil {
		log.Println("❌ Error inserting roles:", err)
	} else {
		log.Println("✔ Roles ready")
	}

	// ===============================
	// 3. INSERT ADMIN
	// ===============================

	adminID := uuid.New().String()
	adminPass, _ := helper.HashPassword("admin123")

	_, err = pg.Exec(`
	INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
	VALUES ($1, 'admin', 'admin@gmail.com', $2, 'Super Admin', $3, true)
	ON CONFLICT (username) DO NOTHING;
	`, adminID, adminPass, roleAdmin)

	if err != nil {
		log.Println("❌ Error inserting admin:", err)
	} else {
		log.Println("✔ Admin user ready")
	}

	// ===============================
	// 4. INSERT LECTURERS (PK Manual)
	// ===============================

	lecturers := []struct {
		ID          string
		UserID      string
		User        string
		Email       string
		LecturerID  string
	}{
		{"LECT-001", uuid.New().String(), "lect1", "lect1@gmail.com", "L001"},
		{"LECT-002", uuid.New().String(), "lect2", "lect2@gmail.com", "L002"},
		{"LECT-003", uuid.New().String(), "lect3", "lect3@gmail.com", "L003"},
	}

	for _, lec := range lecturers {

		pass, _ := helper.HashPassword("lecturer123")

		// Insert user
		_, err = pg.Exec(`
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		ON CONFLICT (username) DO NOTHING;
		`, lec.UserID, lec.User, lec.Email, pass, "Dosen "+lec.User, roleLect)

		if err != nil {
			log.Println("❌ Error inserting lecturer user:", err)
		}

		// Insert lecturer profile
		_, err = pg.Exec(`
		INSERT INTO lecturers (id, user_id, lecturer_id, department)
		VALUES ($1, $2, $3, 'Informatika')
		ON CONFLICT (lecturer_id) DO NOTHING;
		`, lec.ID, lec.UserID, lec.LecturerID)

		if err != nil {
			log.Println("❌ Error inserting lecturer profile:", err)
		}
	}

	log.Println("✔ Lecturers ready")

	// ===============================
	// 5. INSERT STUDENTS
	// ===============================

	students := []struct {
		ID         string
		UserID     string
		User       string
		Email      string
		StudentID  string
		AdvisorPK  string
	}{
		{"STUD-001", uuid.New().String(), "stud1", "stud1@gmail.com", "S001", "LECT-001"},
		{"STUD-002", uuid.New().String(), "stud2", "stud2@gmail.com", "S002", "LECT-002"},
		{"STUD-003", uuid.New().String(), "stud3", "stud3@gmail.com", "S003", "LECT-003"},
	}

	for _, s := range students {

		pass, _ := helper.HashPassword("student123")

		// Insert user
		_, err = pg.Exec(`
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		ON CONFLICT (username) DO NOTHING;
		`, s.UserID, s.User, s.Email, pass, "Mahasiswa "+s.User, roleStud)

		if err != nil {
			log.Println("❌ Error inserting student user:", err)
		}

		// Insert student profile
		_, err = pg.Exec(`
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id)
		VALUES ($1, $2, $3, 'Teknik Informatika', '2023', $4)
		ON CONFLICT (student_id) DO NOTHING;
		`, s.ID, s.UserID, s.StudentID, s.AdvisorPK)

		if err != nil {
			log.Println("❌ Error inserting student profile:", err)
		}
	}

	log.Println("✔ Students ready")
	fmt.Println("🎉 Migration Completed!")
}
