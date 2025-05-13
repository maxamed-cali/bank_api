package controllers

import (
	"bank/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssignRolesInput struct {
	UserID    uint     `json:"user_id" binding:"required"`
	RoleNames []string `json:"roles" binding:"required"`
}



// Create a new role
func CreateRole(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := services.CreateRole(input.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"role": role})
}
func AssignRoles(c *gin.Context) {
	var input AssignRolesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.AssignRolesToUser(input.UserID, input.RoleNames); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Roles assigned successfully"})
}
