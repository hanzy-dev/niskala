package main

import (
	"context"
	"log"
	"time"

	"github.com/hanzy-dev/niskala/apps/api/internal/config"
	"github.com/hanzy-dev/niskala/apps/api/internal/database"
	"github.com/hanzy-dev/niskala/apps/api/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env secara eksplisit di awal
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file tidak ditemukan, menggunakan environment variable sistem")
	}

	cfg := config.Load()

	var dbPool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("Critical: Gagal inisialisasi pool database: %v", err)
		}

		if err := pool.Ping(ctx); err != nil {
			log.Fatalf("Critical: Gagal koneksi ke Postgres (Ping error): %v", err)
		}

		dbPool = pool
		defer dbPool.Close()
		log.Println("Database: Connected successfully to Supabase")
	} else {
		log.Fatal("Critical: DATABASE_URL tidak ditemukan di environment")
	}

	router := server.NewRouter(server.Dependencies{
		DB:                    dbPool,
		PricingServiceBaseURL: cfg.PricingServiceBaseURL,
	})

	log.Printf("Server: Starting on port %s in %s mode", cfg.Port, cfg.AppEnv)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Critical: Gagal menjalankan server: %v", err)
	}
}
