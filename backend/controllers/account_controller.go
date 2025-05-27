package controllers

import (
	"net/http"
	"strconv"

	"bank/models"
	"bank/services"

	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {
	var body models.Account
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.CreateAccount(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, body)
}
func GetAllAccounts(c *gin.Context) {
	accounts, err := services.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func GetAccountsBalance(c *gin.Context) {
	id := c.Param("id")

	accounts, err := services.GetAccountBalance(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func GetAccountDetails(c *gin.Context) {
        accountNumber := c.Query("account_number")
        if accountNumber == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Account number is required"})
            return
        }

        user, err := services.GetUserByAccountNumber(accountNumber)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        c.JSON(http.StatusOK, user)
    }




func GetAccountByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	acc, err := services.GetAccountsByUserID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, acc)
}

func UpdateAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body models.Account
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.UpdateAccount(uint(id), &body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

func DeleteAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeleteAccount(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
