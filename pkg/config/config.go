package config

import (
	"log"
	"os"
)

var (
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
)

func Load() {
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbName = os.Getenv("DB_NAME")
	if DbHost == "" || DbPort == "" || DbUser == "" || DbPassword == "" || DbName == "" {
		log.Fatal("Database environment variables are not set")
	}
}
