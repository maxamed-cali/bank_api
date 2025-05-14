package main

import (
	"bank/db"
	"bank/jobs"
	"bank/routes"
	"bank/websocket"

	"github.com/gin-gonic/gin"
)

func main() {

// Start WebSocket dispatcher
    jobs.StartAutoExpireJob()
    websocket.StartDispatcher()

	r:=gin.Default();
	db.Connect()
	routes.AuthRoutes(r)

	r.Run(":8080") // listen and serve on


}