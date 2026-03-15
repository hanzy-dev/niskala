package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-Debug-User-ID")
		if userID == "" {
			httpx.Unauthorized(c, "Missing authenticated user context")
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
			httpx.Forbidden(c, "Admin access is required")
			c.Abort()
			return
		}

		c.Next()
	}
}
