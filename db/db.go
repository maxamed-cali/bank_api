package db

import (
	"fmt"
	"log"
	"os"

    
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
    "github.com/joho/godotenv"
)

var DB *gorm.DB

func Connect() {

    errs := godotenv.Load()
    if errs != nil {
        log.Fatal("Error loading .env file")
    }
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto-migrate models
    err = DB.AutoMigrate(
        
    )
    if err != nil {
        log.Fatal("Failed to auto migrate:", err)
    }
}


