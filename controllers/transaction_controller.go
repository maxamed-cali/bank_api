package controllers

import (
	"net/http"
	"strconv"

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
