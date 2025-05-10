package db

import (
	"bank/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
        &models.Credential{},
        &models.Role{},
        &models.UserRole{},
        &models.Account{},
        &models.AccountType{},
        &models.Transaction{},
        &models.MoneyRequest{},
        &models.Notification{},
        
    )
    if err != nil {
        log.Fatal("Failed to auto migrate:", err)
    }
}


