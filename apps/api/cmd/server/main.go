package main

import (
	"context"
	"log"

	"github.com/hanzy-dev/niskala/apps/api/internal/config"
	"github.com/hanzy-dev/niskala/apps/api/internal/database"
	"github.com/hanzy-dev/niskala/apps/api/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	var dbPool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		pool, err := database.NewPostgresPool(context.Background(), cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("failed to connect postgres: %v", err)
		}

		dbPool = pool
		defer dbPool.Close()
	}

	router := server.NewRouter(server.Dependencies{
		DB:                    dbPool,
		PricingServiceBaseURL: cfg.PricingServiceBaseURL,
	})

	log.Printf("starting api server on port %s in %s mode", cfg.Port, cfg.AppEnv)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start api server: %v", err)
	}
}
