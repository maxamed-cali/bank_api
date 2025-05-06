package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"bank/models"
	"bank/services"
)

// Create wallet handler
func CreateWallet(c *gin.Context) {
	var wallet models.Account
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.CreateWallet(&wallet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, wallet)
}

// Delete wallet handler
func DeleteWallet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeleteWallet(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted"})
}

// Rename wallet
func RenameWallet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body struct {
		NewAccountNumber string `json:"new_account_number"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.RenameWallet(uint(id), body.NewAccountNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Wallet renamed"})
}


func ViewBalances(c *gin.Context) {
	accountNumber := c.Param("accountNumber") 
	balance, err := services.GetBalanceByAccountNumber(accountNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"account_number": accountNumber,
		"balance":        balance,
	})
}


// View balances by currency
func HandleGetBalance(c *gin.Context) {
	result, err := services.GetBalancesGroupedByCurrency()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
