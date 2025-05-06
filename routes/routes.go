package routes

import (
	"bank/controllers"
	"bank/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	api := r.Group("/api")
	api.Use(middlewares.JWTAuthMiddleware()) // Only authenticated


	{
		
		user := api.Group("/user")
		user.Use(middlewares.RoleMiddleware("User"))
		{
			user.GET("/dashboard", controllers.UserDashboard)
		}

		admin := api.Group("/admin")
		admin.Use(middlewares.RoleMiddleware("Admin"))
		{
			admin.GET("/dashboard", controllers.AdminDashboard)
		}
	}

	
}
