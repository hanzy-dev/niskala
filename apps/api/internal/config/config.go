package config

import "os"

type Config struct {
	AppEnv string
	Port   string
}

func Load() Config {
	appEnv := getEnv("APP_ENV", "development")
	port := getEnv("PORT", "8080")

	return Config{
		AppEnv: appEnv,
		Port:   port,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}