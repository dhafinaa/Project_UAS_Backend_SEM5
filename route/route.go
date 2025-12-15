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
	achRepo := repository.NewAchievementRepository(mongoDb, pg)
	lecturerRepo := repository.NewLecturerRepository(pg)


	// -------------------------------------------
	// SERVICES
	// -------------------------------------------
	authService := service.NewAuthService(authRepo)         
	achievementService := service.NewAchievementService(achRepo, studentRepo) 
	lecturerService := service.NewLecturerService(studentRepo, achRepo, lecturerRepo)

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
	student := app.Group("/student", middleware.AuthRequired(authRepo), middleware.RoleRequired("Mahasiswa"),)

	// ---- ACHIEVEMENTS ----
	student.Get("/achievements", middleware.PermissionRequired("achievement.read"), achievementService.GetAchievements)
	student.Get("/achievements/:id", middleware.PermissionRequired("achievement.read"), achievementService.GetAchievementDetail)
	student.Post("/achievements", middleware.PermissionRequired("achievement.create"), achievementService.CreateAchievement)
	student.Put("/achievements/:id", middleware.PermissionRequired("achievement.update"), achievementService.UpdateAchievement)
	student.Post("/achievements/:id/submit", middleware.PermissionRequired("achievement.submit"), achievementService.SubmitAchievement)
	student.Delete("/achievements/:id",middleware.PermissionRequired("achievement.delete"), achievementService.DeleteAchievement)
	student.Post("/achievements/:id/attachments", middleware.PermissionRequired("achievement.attachment.upload"), achievementService.UploadAttachment)


	// -------------------------------------------
	// LECTURER ROUTES
	// -----------------------------------------x--
	lecturer := app.Group("/lecturer",
		middleware.AuthRequired(authRepo),
		middleware.RoleRequired("Dosen Wali"),
	)

	lecturer.Get("/advisees", middleware.PermissionRequired("advisee.read"), lecturerService.GetStudentAchievements)
	lecturer.Put("/achievements/:id/verify", middleware.PermissionRequired("achievement.verify"), lecturerService.VerifyAchievement)
	lecturer.Put("/achievements/:id/reject", middleware.PermissionRequired("achievement.reject"), lecturerService.RejectAchievement)




}