package route

import (
	"github.com/gofiber/fiber/v2"

	"PROJECT_UAS/app/repository"
	"PROJECT_UAS/app/service"
	"database/sql"
)

func RegisterRoutes(app *fiber.App, pg *sql.DB) {

	// Init Auth Repository
	authRepo := repository.NewAuthRepository(pg)
	authService := service.NewAuthService(authRepo)

	// AUTH ROUTES
	app.Post("/auth/login", authService.Login)
	app.Post("/auth/refresh", authService.RefreshToken)
	app.Post("/auth/logout", authService.Logout)
	app.Get("/auth/profile", authService.Profile)

	/*
	// ============================
	// 4.2 Manajemen Prestasi (Mahasiswa)
	// ============================
	app.Post("/achievements")             // FR-003 Submit Prestasi
	app.Put("/achievements/:id/submit")   // FR-004 Submit verifikasi
	app.Delete("/achievements/:id")       // FR-005 Hapus Prestasi (soft delete)

	// ============================
	// 4.3 Verifikasi Prestasi (Dosen Wali)
	// ============================
	app.Get("/advisor/achievements")      // FR-006 View Prestasi Bimbingan
	app.Put("/achievements/:id/verify")   // FR-007 Verify Prestasi
	app.Put("/achievements/:id/reject")   // FR-008 Reject Prestasi

	// ============================
	// 4.4 Manajemen Sistem (Admin)
	// ============================
	app.Post("/admin/users")               // FR-009 Create user
	app.Get("/admin/users")                // FR-009 List users
	app.Put("/admin/users/:id")            // FR-009 Update user
	app.Delete("/admin/users/:id")         // FR-009 Delete user
	app.Put("/admin/users/:id/role")       // FR-009 Assign role
	app.Put("/admin/students/:id/advisor") // FR-009 Set advisor

	// ============================
	// 4.5 Reporting & Analytics
	// ============================
	app.Get("/reports/achievements")       // FR-011 Statistics
	*/
}
