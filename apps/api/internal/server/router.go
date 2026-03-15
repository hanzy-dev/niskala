package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/auth"
	"github.com/hanzy-dev/niskala/apps/api/internal/handler"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	meHandler := handler.NewMeHandler()

	router.GET("/health", healthHandler.GetHealth)

	api := router.Group("/api")
	{
		protected := api.Group("")
		protected.Use(auth.RequireAuth())
		{
			protected.GET("/me", meHandler.GetMe)

			admin := protected.Group("/admin")
			admin.Use(auth.RequireAdmin())
			{
				admin.GET("/ping", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"status": "ok",
						"scope":  "admin",
					})
				})
			}
		}
	}

	return router
}
