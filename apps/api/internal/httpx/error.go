package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func JSONError(c *gin.Context, status int, code string, message string, details any) {
	c.JSON(status, gin.H{
		"error": ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
		"correlation_id": GetCorrelationID(c),
	})
}

func Unauthorized(c *gin.Context, message string) {
	JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func Forbidden(c *gin.Context, message string) {
	JSONError(c, http.StatusForbidden, "FORBIDDEN", message, nil)
}

func BadRequest(c *gin.Context, code string, message string) {
	JSONError(c, http.StatusBadRequest, code, message, nil)
}

func NotFound(c *gin.Context, code string, message string) {
	JSONError(c, http.StatusNotFound, code, message, nil)
}

func Internal(c *gin.Context, code string, message string) {
	JSONError(c, http.StatusInternalServerError, code, message, nil)
}
