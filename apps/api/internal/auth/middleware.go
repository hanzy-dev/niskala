package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/authjwt"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
)

type Middleware struct {
	verifier *authjwt.Verifier
}

func NewMiddleware(verifier *authjwt.Verifier) *Middleware {
	return &Middleware{
		verifier: verifier,
	}
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))

		if authHeader != "" && m.verifier != nil && m.verifier.Enabled() {
			claims, err := m.verifier.ParseBearerToken(authHeader)
			if err == nil {
				role := claims.Role
				if role == "" {
					role = RoleUser
				}

				SetUserContext(c, claims.Subject, role)
				c.Next()
				return
			}
		}

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

func (m *Middleware) RequireAdmin() gin.HandlerFunc {
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
