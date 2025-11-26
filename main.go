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

	godotenv.Load()

	db := database.ConnectMongo()
	_ = db // sementara tidak dipakai, tapi harus tetap connect

	app := fiber.New()

	// Semua endpoint digabung jadi satu
	route.RegisterRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running at http://localhost:" + port)
	app.Listen(":" + port)
}