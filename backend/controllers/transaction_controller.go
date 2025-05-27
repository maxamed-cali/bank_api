package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"bank/models"
	"bank/services"

	"github.com/gin-gonic/gin"
)

func MoneyTransfer(c *gin.Context) {
	var tx models.Transaction

	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.MoneyTransfer(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transfer successful",
		"transaction": tx,
	})
}

func MoneyRequest(c *gin.Context) {
	var mr models.MoneyRequest

	if err := c.ShouldBindJSON(&mr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.MoneyRequest(&mr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Money Requst successful",
		"Money":   mr,
	})
}

func AcceptMoneyRequest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := services.AcceptMoneyRequest(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Accepted"})
}

func DeclineMoneyRequest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := services.DeclineMoneyRequest(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Declined"})
}

func GetTransactionHistoryHandler(c *gin.Context) {
	var filter services.TransactionFilter

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			uid := uint(userID)
			filter.UserID = &uid
		}
	}
	if accountID := c.Query("account_id"); accountID != "" {
		filter.AccountID = &accountID
	}
	if transactionType := c.Query("transaction_type"); transactionType != "" {
		filter.TransactionType = &transactionType
	}
	if minAmountStr := c.Query("min_amount"); minAmountStr != "" {
		if minAmount, err := strconv.ParseFloat(minAmountStr, 64); err == nil {
			filter.MinAmount = &minAmount
		}
	}
	if maxAmountStr := c.Query("max_amount"); maxAmountStr != "" {
		if maxAmount, err := strconv.ParseFloat(maxAmountStr, 64); err == nil {
			filter.MaxAmount = &maxAmount
		}
	}
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}
	if description := c.Query("description"); description != "" {
		filter.DescriptionLike = &description
	}

	transactions, err := services.GetTransactionHistory(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

func GetMoneyRequestsByUserID(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	requests, err := services.GetMoneyRequestsByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func GetFilteredNotifications(c *gin.Context) {
	userID := c.Query("user_id")
	filter := c.Query("filter") // all, requests, alert

	if userID == "" {
		c.JSON(400, gin.H{"error": "user_id is required"})
		return
	}

	notifications, err := services.GetFilteredNotifications(userID, filter)
	if err != nil {
		log.Printf("Error fetching notifications: %v", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve notifications"})
		return
	}

	c.JSON(200, notifications)
}

// GetDashboard returns summary data for the user's dashboard
func GetDashboard(c *gin.Context) {
	// Extract userID from context (assumes middleware sets it)
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Call service to get dashboard summary
	summary, err := services.GetDashboardSummary(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the dashboard summary
	c.JSON(http.StatusOK, gin.H{
		"data": summary,
	})
}

func GetMonthlyTransactionVolume(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	data, err := services.GetMonthlyTransactionVolume(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
