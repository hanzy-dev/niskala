package httpx

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const CorrelationIDHeader = "X-Correlation-ID"
const CorrelationIDContextKey = "http.correlation_id"

func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = newCorrelationID()
		}

		c.Set(CorrelationIDContextKey, correlationID)
		c.Writer.Header().Set(CorrelationIDHeader, correlationID)
		c.Next()
	}
}

func GetCorrelationID(c *gin.Context) string {
	value, exists := c.Get(CorrelationIDContextKey)
	if !exists {
		return ""
	}

	correlationID, ok := value.(string)
	if !ok {
		return ""
	}

	return correlationID
}

func newCorrelationID() string {
	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		return "generated-correlation-id"
	}

	return hex.EncodeToString(buffer)
}
