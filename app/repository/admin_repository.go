package repository

// import (
// 	"database/sql"
// 	"errors"
// 	"PROJECT_UAS/app/model"
// 	"PROJECT_UAS/helper"
// )

// type AdminRepository struct {
// 	DB *sql.DB
// }

// func NewAdminRepository(db *sql.DB) *AdminRepository {
// 	return &AdminRepository{DB: db}
// }

// /* ======================
//    CREATE USER
// ====================== */
// func (r *AdminRepository) CreateUser(req model.CreateUserRequest) error {

// 	hash, _ := helper.HashPassword(req.Password)

// 	query := `
// 		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active)
// 		VALUES ($1, $2, $3, $4, $5, true)
// 	`

// 	_, err := r.DB.Exec(query,
// 		req.Username,
// 		req.Email,
// 		hash,
// 		req.FullName,
// 		req.RoleID,
// 	)

// 	return err
// }

// /* ======================
//    UPDATE USER
// ====================== */
// func (r *AdminRepository) UpdateUser(id string, req model.UpdateUserRequest) error {

// 	query := `
// 		UPDATE users 
// 		SET full_name = $1, email = $2, is_active = $3 
// 		WHERE id = $4
// 	`

// 	_, err := r.DB.Exec(query, req.FullName, req.Email, req.IsActive, id)
// 	return err
// }

// /* ======================
//    DELETE USER
// ====================== */
// func (r *AdminRepository) DeleteUser(id string) error {
// 	_, err := r.DB.Exec(`DELETE FROM users WHERE id=$1`, id)
// 	return err
// }

// /* ======================
//    ASSIGN ROLE
// ====================== */
// func (r *AdminRepository) AssignRole(id string, roleID string) error {

// 	_, err := r.DB.Exec(`UPDATE users SET role_id=$1 WHERE id=$2`, roleID, id)
// 	return err
// }

// /* ===============================
//    SET STUDENT PROFILE
// ================================*/
// func (r *AdminRepository) SetStudentProfile(userID string, s model.Student) error {

// 	query := `
// 		INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id)
// 		VALUES ($1, $2, $3, $4, $5)
// 		ON CONFLICT (user_id)
// 		DO UPDATE 
// 		   SET student_id = EXCLUDED.student_id,
// 		       program_study = EXCLUDED.program_study,
// 		       academic_year = EXCLUDED.academic_year,
// 		       advisor_id = EXCLUDED.advisor_id
// 	`
// 	_, err := r.DB.Exec(query,
// 		userID,
// 		s.Student_id,
// 		s.Program_study,
// 		s.Academic_year,
// 		s.Advisor_id,
// 	)

// 	return err
// }

// /* ===============================
//    SET LECTURER PROFILE
// ================================*/
// func (r *AdminRepository) SetLecturerProfile(userID string, l model.Lecturer) error {

// 	query := `
// 		INSERT INTO lecturers (user_id, lecturer_id, department)
// 		VALUES ($1, $2, $3)
// 		ON CONFLICT (user_id)
// 		DO UPDATE 
// 		   SET lecturer_id = EXCLUDED.lecturer_id,
// 		       department = EXCLUDED.department
// 	`
// 	_, err := r.DB.Exec(query,
// 		userID,
// 		l.Lecturer_id,
// 		l.Department,
// 	)

// 	return err
// }

// /* ===============================
//    SET ADVISOR FOR STUDENT
// ================================*/
// func (r *AdminRepository) SetAdvisor(studentUserID string, advisorID string) error {

// 	query := `
// 		UPDATE students 
// 		SET advisor_id = $1
// 		WHERE user_id = $2
// 	`

// 	_, err := r.DB.Exec(query, advisorID, studentUserID)
// 	return err
// }