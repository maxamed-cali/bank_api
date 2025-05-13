package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"bank/models"
	"bank/services"
)

func CreateAccountType(c *gin.Context) {
	var input models.AccountType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.CreateAccountType(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func GetAllAccountTypes(c *gin.Context) {
	list, err := services.GetAllAccountTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func GetAccountTypeByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	result, err := services.GetAccountTypeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func UpdateAccountType(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body models.AccountType
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.UpdateAccountType(uint(id), &body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

func DeleteAccountType(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeleteAccountType(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
