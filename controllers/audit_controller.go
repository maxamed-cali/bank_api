package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bank/services"
)

// GET /audit-logs
func GetAuditLogs(c *gin.Context) {
	logs, err := services.GetAllAuditLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch audit logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
	})
}
