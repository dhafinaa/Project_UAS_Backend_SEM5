package route

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"PROJECT_UAS/app/repository"
	"PROJECT_UAS/app/service"
	"PROJECT_UAS/middleware"
)

func RegisterRoutes(app *fiber.App, pg *sql.DB, mongoDb *mongo.Database) {

	// -------------------------------------------
	// REPOSITORIES
	// -------------------------------------------
	authRepo := repository.NewAuthRepository(pg)
	studentRepo := &repository.StudentRepository{DB: pg}
	achRepo := repository.NewAchievementRepository(mongoDb)

	// -------------------------------------------
	// SERVICES
	// -------------------------------------------
	authService := service.NewAuthService(authRepo)
	studentService := service.NewStudentService(achRepo, studentRepo)
	lecturerService := service.NewLecturerService(studentRepo, achRepo)

	// -------------------------------------------
	// AUTH ROUTES
	// -------------------------------------------
	auth := app.Group("/auth")
	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.RefreshToken)
	auth.Post("/logout", authService.Logout)
	auth.Get("/profile", middleware.AuthRequired(authRepo), authService.Profile)

	// -------------------------------------------
	// STUDENT ROUTES
	// -------------------------------------------
	student := app.Group("/student",
		middleware.AuthRequired(authRepo),
		middleware.RoleRequired("Mahasiswa"),
	)

	student.Get("/achievements", middleware.PermissionRequired("achievement.read"), studentService.GetAchievements)
	student.Post("/achievements", middleware.PermissionRequired("achievement.create"), studentService.CreateAchievement)
	student.Put("/achievements/:id/submit", middleware.PermissionRequired("achievement.submit"), studentService.SubmitAchievement)
	student.Delete("/achievements/:id", middleware.PermissionRequired("achievement.delete"), studentService.DeleteAchievement)

	// -------------------------------------------
	// LECTURER ROUTES
	// -------------------------------------------
	lecturer := app.Group("/lecturer",
		middleware.AuthRequired(authRepo),
		middleware.RoleRequired("Dosen Wali"),
	)

	lecturer.Get("/advisees", middleware.PermissionRequired("advisee.read"), lecturerService.GetStudentAchievements)
	lecturer.Put("/achievements/:id/verify", middleware.PermissionRequired("achievement.verify"), lecturerService.VerifyAchievement)
	lecturer.Put("/achievements/:id/reject", middleware.PermissionRequired("achievement.reject"), lecturerService.RejectAchievement)



	// ADMIN
	// admin := app.Group("/admin",
	// 	middleware.AuthRequired(authRepo),
	// 	middleware.RoleRequired("Admin"),
	// )
	// admin.Post("/users", middleware.PermissionRequired("user.manage"), adminService.CreateUser)
	// admin.Put("/users/:id", middleware.PermissionRequired("user.manage"), adminService.UpdateUser)
	// admin.Delete("/users/:id", middleware.PermissionRequired("user.manage"), adminService.DeleteUser)
	// admin.Put("/users/:id/role", middleware.PermissionRequired("user.manage"), adminService.AssignRole)

	// admin.Put("/students/:id/profile", middleware.PermissionRequired("user.manage"), adminService.SetStudentProfile)
	// admin.Put("/lecturers/:id/profile", middleware.PermissionRequired("user.manage"), adminService.SetLecturerProfile)
	// admin.Put("/students/:id/advisor", middleware.PermissionRequired("student.assign.advisor"), adminService.SetAdvisor)

	// admin.Get("/reports/achievements", middleware.PermissionRequired("reports.read"), adminService.GenerateAchievementReport)
}