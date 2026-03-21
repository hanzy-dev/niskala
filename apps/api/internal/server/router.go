package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hanzy-dev/niskala/apps/api/internal/auth"
	"github.com/hanzy-dev/niskala/apps/api/internal/authjwt"
	"github.com/hanzy-dev/niskala/apps/api/internal/handler"
	"github.com/hanzy-dev/niskala/apps/api/internal/httpx"
	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
	"github.com/hanzy-dev/niskala/apps/api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	DB                    *pgxpool.Pool
	Redis                 *redis.Client
	PricingServiceBaseURL string
	JWTVerifier           *authjwt.Verifier
}

func NewRouter(deps Dependencies) *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Correlation-ID", "Idempotency-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(httpx.CorrelationIDMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	productRepository := repository.NewProductRepository(deps.DB)
	cartRepository := repository.NewCartRepository(deps.DB)
	orderRepository := repository.NewOrderRepository(deps.DB)
	checkoutRepository := repository.NewCheckoutRepository(deps.DB)
	idempotencyRepository := repository.NewIdempotencyRepository(deps.DB)
	membershipRepository := repository.NewMembershipRepository(deps.DB)

	productService := service.NewProductService(productRepository)
	productCacheService := service.NewProductCacheService(deps.Redis)
	cartService := service.NewCartService(cartRepository)
	orderService := service.NewOrderService(orderRepository, checkoutRepository)
	idempotencyService := service.NewIdempotencyService(idempotencyRepository, orderRepository)
	pricingService := service.NewPricingService(deps.PricingServiceBaseURL)
	membershipService := service.NewMembershipService(membershipRepository)

	checkoutService := service.NewCheckoutService(
		productService,
		cartService,
		orderService,
		idempotencyService,
		pricingService,
	)

	authMiddleware := auth.NewMiddleware(deps.JWTVerifier, membershipService)

	healthService := service.NewHealthService(deps.DB, deps.PricingServiceBaseURL)
	healthHandler := handler.NewHealthHandler(healthService)

	productHandler := handler.NewProductHandler(productService, productCacheService)
	adminProductHandler := handler.NewAdminProductHandler(productService, productCacheService)
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
		protected.Use(authMiddleware.RequireAuth())
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
			admin.Use(authMiddleware.RequireAdmin())
			{
				admin.GET("/ping", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"status": "ok",
						"scope":  "admin",
					})
				})

				admin.POST("/products", adminProductHandler.CreateProduct)
				admin.PATCH("/products/:id", adminProductHandler.UpdateProduct)
				admin.PATCH("/products/:id/stock", adminProductHandler.UpdateStock)
			}
		}
	}

	return router
}
