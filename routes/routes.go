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

			user.POST("/account", controllers.CreateWallet)
			user.GET("/account/:accountNumber", controllers.ViewBalances)
			user.PUT("/account/:id", controllers.RenameWallet)
			user.DELETE("/account/:id", controllers.DeleteWallet)
			user.GET("/account", controllers.HandleGetBalance)

			user.POST("/account-types", controllers.CreateAccountType)
			user.GET("/account-types", controllers.GetAllAccountTypes)
			user.GET("/account-types/:id", controllers.GetAccountTypeByID)
			user.PUT("/account-types/:id", controllers.UpdateAccountType)
			user.DELETE("/account-types/:id", controllers.DeleteAccountType)
			user.GET("/dashboard", controllers.UserDashboard)
		}

		admin := api.Group("/admin")
		admin.Use(middlewares.RoleMiddleware("Admin"))
		{
			admin.GET("/dashboard", controllers.AdminDashboard)

		}
	}

}
