package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    _"github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func Connect() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Format DSN
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    // Open database connection
    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal("Failed to open database connection:", err)
    }

    // Test the connection
    err = DB.Ping()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    fmt.Println("Database connection established using raw SQL")
}

func GetDB() *sql.DB {
	if DB == nil {
		log.Fatal("Database not connected. Call db.Connect() first.")
	}
	return DB
}