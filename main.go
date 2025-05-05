package main

import (
	"bank/db"
	"bank/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	r:=gin.Default();
	db.Connect()
	routes.AuthRoutes(r)

	r.Run(":8080") // listen and serve on


}