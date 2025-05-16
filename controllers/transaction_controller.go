package controllers

import (
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
		"message":     "Money Requst successful",
		"Money": mr,
	})
}

func AcceptMoneyRequest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	 err := services.AcceptMoneyRequest(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
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