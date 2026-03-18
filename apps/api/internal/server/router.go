package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/auth"
	"github.com/hanzy-dev/niskala/apps/api/internal/handler"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	DB                    *pgxpool.Pool
	PricingServiceBaseURL string
}

func NewRouter(deps Dependencies) *gin.Engine {
	router := gin.New()

	router.Use(httpx.CorrelationIDMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	healthService := service.NewHealthService(deps.DB, deps.PricingServiceBaseURL)
	healthHandler := handler.NewHealthHandler(healthService)

	productRepository := repository.NewProductRepository(deps.DB)
	cartRepository := repository.NewCartRepository(deps.DB)

	productService := service.NewProductService(productRepository)
	cartService := service.NewCartService(cartRepository)
	orderService := service.NewOrderService()
	idempotencyService := service.NewIdempotencyService()
	pricingService := service.NewPricingService(deps.PricingServiceBaseURL)
	checkoutService := service.NewCheckoutService(
		productService,
		cartService,
		orderService,
		idempotencyService,
		pricingService,
	)

	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	checkoutHandler := handler.NewCheckoutHandler(checkoutService)
	orderHandler := handler.NewOrderHandler(orderService)
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

			protected.POST("/checkout", checkoutHandler.Checkout)
			protected.GET("/orders", orderHandler.ListOrders)
			protected.GET("/orders/:id", orderHandler.GetOrder)

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
