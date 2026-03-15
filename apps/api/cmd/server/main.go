package main

import (
	"log"

	"github.com/hanzy-dev/niskala/apps/api/internal/config"
	"github.com/hanzy-dev/niskala/apps/api/internal/server"
)

func main() {
	cfg := config.Load()
	router := server.NewRouter()

	log.Printf("starting api server on port %s in %s mode", cfg.Port, cfg.AppEnv)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start api server: %v", err)
	}
}