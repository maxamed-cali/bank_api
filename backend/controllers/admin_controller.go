package controllers

import (
	"bank/services"
	"net/http"
	"strconv"

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

type StatusToggleRequest struct {
	IsActive bool `json:"is_active"`
}
func ActivateDeactivateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req StatusToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	err = services.ActivateDeactivateUser(uint(userID), req.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusMsg := "deactivated"
	if req.IsActive {
		statusMsg = "activated"
	}
	c.JSON(http.StatusOK, gin.H{"message": "User " + statusMsg + " successfully"})
}

func GetUserList(c *gin.Context) {
  
    users, err := services.GetUsersWithRoles()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}
