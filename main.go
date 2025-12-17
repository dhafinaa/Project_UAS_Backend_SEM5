package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "PROJECT_UAS/docs" 

	"PROJECT_UAS/database"
	"PROJECT_UAS/route"
	"PROJECT_UAS/middleware"
)

func main() {

	// Load ENV
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system ENV")
	}

	// ==============================
	// CONNECT MONGODB
	// ==============================
	mongoDB := database.ConnectMongo()

	// ==============================
	// CONNECT POSTGRES
	// ==============================
	postgresDB := database.ConnectPostgres()

	blacklist := middleware.NewTokenBlacklist()

	// ==============================
	// START FIBER APP
	// ==============================
	app := fiber.New()

	// ==============================
	// SWAGGER ROUTE (WAJIB BIAR GA ERROR)
	// ==============================
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// ==============================
	// REGISTER ROUTES
	// ==============================
	route.RegisterRoutes(app, postgresDB, mongoDB, blacklist)

	// ==============================
	// RUN SERVER
	// ==============================
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running at http://localhost:" + port)
	log.Println("Swagger running at http://localhost:" + port + "/swagger/index.html")

	log.Fatal(app.Listen(":" + port))
}