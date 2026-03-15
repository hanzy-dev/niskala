package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/handler"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()

	router.GET("/health", healthHandler.GetHealth)

	return router
}