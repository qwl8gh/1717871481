package main

import (
	"log"
	"web-messaging-service/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	config.Load()

	// Database connection string
	dbConnStr := "postgres://" + config.DbUser + ":" + config.DbPassword + "@" + config.DbHost + ":" + config.DbPort + "/" + config.DbName + "?sslmode=disable"

	// Create a new migration instance
	m, err := migrate.New(
		"file://./migrations",
		dbConnStr)
	if err != nil {
		log.Fatalf("could not create migration instance: %v", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("could not run migrations: %v", err)
	}

	log.Println("Database migration completed successfully")
}
