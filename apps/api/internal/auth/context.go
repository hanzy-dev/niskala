package auth

import "github.com/gin-gonic/gin"

const (
	ContextUserIDKey   = "auth.user_id"
	ContextUserRoleKey = "auth.user_role"
)

func SetUserContext(c *gin.Context, userID string, role string) {
	c.Set(ContextUserIDKey, userID)
	c.Set(ContextUserRoleKey, role)
}

func GetUserID(c *gin.Context) (string, bool) {
	value, exists := c.Get(ContextUserIDKey)
	if !exists {
		return "", false
	}

	userID, ok := value.(string)
	return userID, ok
}

func GetUserRole(c *gin.Context) (string, bool) {
	value, exists := c.Get(ContextUserRoleKey)
	if !exists {
		return "", false
	}

	role, ok := value.(string)
	return role, ok
}