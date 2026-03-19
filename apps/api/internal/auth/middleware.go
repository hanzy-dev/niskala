package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/authjwt"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
)

type Middleware struct {
	verifier          *authjwt.Verifier
	membershipService *service.MembershipService
}

func NewMiddleware(verifier *authjwt.Verifier, membershipService *service.MembershipService) *Middleware {
	return &Middleware{
		verifier:          verifier,
		membershipService: membershipService,
	}
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))

		if authHeader != "" && m.verifier != nil && m.verifier.Enabled() {
			claims, err := m.verifier.ParseBearerToken(authHeader)
			if err == nil {
				SetUserContext(c, claims.Subject, RoleUser)
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

		SetUserContext(c, userID, RoleUser)
		c.Next()
	}
}

func (m *Middleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok || userID == "" {
			httpx.Unauthorized(c, "Missing authenticated user context")
			c.Abort()
			return
		}

		isAdmin, err := m.membershipService.IsAdmin(c.Request.Context(), userID)
		if err != nil {
			httpx.Internal(c, "AUTHORIZATION_FAILED", "Failed to verify admin membership")
			c.Abort()
			return
		}

		if !isAdmin {
			httpx.Forbidden(c, "Admin access is required")
			c.Abort()
			return
		}

		c.Next()
	}
}
