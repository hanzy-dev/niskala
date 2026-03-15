package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/auth"
	"github.com/hanzy-dev/niskala/apps/api/internal/handler"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	productService := service.NewProductService()
	cartService := service.NewCartService()

	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	meHandler := handler.NewMeHandler()

	router.GET("/health", healthHandler.GetHealth)

	api := router.Group("/api")
	{
		api.GET("/products", productHandler.ListProducts)
		api.GET("/products/:id", productHandler.GetProduct)

		protected := api.Group("")
		protected.Use(auth.RequireAuth())
		{
			protected.GET("/me", meHandler.GetMe)
			protected.GET("/cart", cartHandler.GetCart)
			protected.POST("/cart/items", cartHandler.AddCartItem)
			protected.PATCH("/cart/items/:productId", cartHandler.UpdateCartItem)
			protected.DELETE("/cart/items/:productId", cartHandler.DeleteCartItem)

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
