package config

import "os"

type Config struct {
	AppEnv                string
	Port                  string
	DatabaseURL           string
	PricingServiceBaseURL string
}

func Load() Config {
	appEnv := getEnv("APP_ENV", "development")
	port := getEnv("PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "")
	pricingServiceBaseURL := getEnv("PRICING_SERVICE_BASE_URL", "http://localhost:8081")

	return Config{
		AppEnv:                appEnv,
		Port:                  port,
		DatabaseURL:           databaseURL,
		PricingServiceBaseURL: pricingServiceBaseURL,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
