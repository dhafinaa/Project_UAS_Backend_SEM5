package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"PROJECT_UAS/database"
	"PROJECT_UAS/route"
)

func main() {

	// Load ENV
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system ENV")
	}

	// ==============================
	// CONNECT MONGODB (required)
	// ==============================
	mongoDB := database.ConnectMongo()
	_ = mongoDB // nanti dipakai di achievement repo

	// ==============================
	// CONNECT POSTGRES (required)
	// ==============================
	postgresDB := database.ConnectPostgres()

	// ==============================
	// RUN MIGRATION (mengisi users, roles, dll)
	// ==============================
	// database.RunMigration(postgresDB)

	// ==============================
	// START FIBER APP
	// ==============================
	app := fiber.New()

	// kirim postgres ke route
	route.RegisterRoutes(app, postgresDB, mongoDB)

	// ==============================
	// RUN SERVER
	// ==============================
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running at http://localhost:" + port)
	app.Listen(":" + port)
}
