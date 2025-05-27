package main

import (
	"bank/db"
	"bank/jobs"
	"bank/routes"
	"bank/websocket"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database Connection
	db.Connect()

	// Run Migrations Before Starting Server
	if err := db.RunMigrations(db.GetDB()); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Start Background Jobs and WebSocket Dispatcher
	jobs.StartAutoExpireJob()
	websocket.StartDispatcher()

	// Set up Gin Router
	r := gin.Default()

	// CORS Configuration
	config := cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"*"},
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
	}

	r.Use(cors.New(config))

	// Register Routes
	routes.AuthRoutes(r)

	// Start HTTP Server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
