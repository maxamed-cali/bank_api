package routes

import (
	"bank/controllers"
	"bank/middlewares"
	"bank/websocket"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
    r.GET("/ws", websocket.WebSocketHandler)
	api := r.Group("/api")
	api.Use(middlewares.JWTAuthMiddleware()) // Only authenticated

	{

		user := api.Group("/user")
		user.Use(middlewares.RoleMiddleware("User"))
		{
 
			user.POST("/password-reset", controllers.ResetPassword)
			user.POST("/accounts", controllers.CreateAccount)
			user.PUT("/accounts/:id", controllers.UpdateAccount)
			user.GET("/accounts", controllers.GetAllAccounts)
			user.GET("/accounts/:id", controllers.GetAccountByID)
			user.GET("/account-blance/:id", controllers.GetAccountsBalance)
			user.DELETE("/accounts/:id", controllers.DeleteAccount)
			
			user.POST("/account-types", controllers.CreateAccountType)
			user.GET("/account-types", controllers.GetAllAccountTypes)
			user.GET("/account-types/:id", controllers.GetAccountTypeByID)
			user.PUT("/account-types/:id", controllers.UpdateAccountType)
			user.DELETE("/account-types/:id", controllers.DeleteAccountType)
			
         	user.POST("/money-transer", controllers.MoneyTransfer)
			user.GET("/transactions/money-request", controllers.GetMoneyRequestsByUserID)
			user.GET("/transactions/history", controllers.GetTransactionHistoryHandler)
			user.POST("/money-request", controllers.MoneyRequest)
			user.PUT("/accept-money-request/:id", controllers.AcceptMoneyRequest)
			user.PUT("/decline-money-request/:id", controllers.DeclineMoneyRequest)
			user.GET("/notifications",controllers.GetFilteredNotifications)
			user.GET("/dashboard/transactions-summary", controllers.GetDashboard)
			user.GET("/dashboard/monthly-transactions", controllers.GetMonthlyTransactionVolume)
			user.GET("/account-details", controllers.GetAccountDetails)


		}

		admin := api.Group("/admin")
		admin.Use(middlewares.RoleMiddleware("Admin"))
		{
			
			admin.POST("/assign-roles", controllers.AssignRoles)
			admin.PUT("/users/:id/status", controllers.ActivateDeactivateUser)
			admin.POST("/create-role", controllers.CreateRole)
			admin.GET("/users", controllers.GetUserList)
			admin.GET("/audit-logs", controllers.GetAuditLogs)

			
		}
	}

}
