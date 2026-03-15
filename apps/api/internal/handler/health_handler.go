package handler

import "github.com/gin-gonic/gin"

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":             "ok",
		"service":            "api",
		"database":           "unknown",
		"pricing_service":    "unknown",
		"checkout_available": false,
	})
}
