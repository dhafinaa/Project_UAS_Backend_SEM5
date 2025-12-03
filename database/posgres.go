package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectPostgres() *sql.DB {

	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DB")

	if host == "" || port == "" || user == "" || pass == "" || dbname == "" {
		log.Fatal("❌ PostgreSQL environment variables are missing")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, dbname, port,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Failed to open PostgreSQL connection:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("❌ Cannot connect to PostgreSQL:", err)
	}

	log.Println("🐘 Connected to PostgreSQL")
	return db
}