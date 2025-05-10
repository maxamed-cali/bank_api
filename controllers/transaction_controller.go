package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bank/models"
	"bank/services"
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
