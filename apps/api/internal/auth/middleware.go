package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-Debug-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Missing authenticated user context",
					"details": nil,
				},
			})
			c.Abort()
			return
		}

		role := c.GetHeader("X-Debug-User-Role")
		if role == "" {
			role = RoleUser
		}

		SetUserContext(c, userID, role)
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetUserRole(c)
		if !ok || !IsAdmin(role) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Admin access is required",
					"details": nil,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
