package config

import "os"

type Config struct {
	AppEnv                string
	Port                  string
	DatabaseURL           string
	PricingServiceBaseURL string
	SupabaseURL           string
	SupabaseJWKSURL       string
	SupabaseJWTAudience   string
	SupabaseJWTIssuer     string
}

func Load() Config {
	appEnv := getEnv("APP_ENV", "development")
	port := getEnv("PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "")
	pricingServiceBaseURL := getEnv("PRICING_SERVICE_BASE_URL", "http://localhost:8081")

	supabaseURL := getEnv("SUPABASE_URL", "")
	supabaseJWKSURL := getEnv("SUPABASE_JWKS_URL", "")
	supabaseJWTAudience := getEnv("SUPABASE_JWT_AUDIENCE", "authenticated")
	supabaseJWTIssuer := getEnv("SUPABASE_JWT_ISSUER", "")

	return Config{
		AppEnv:                appEnv,
		Port:                  port,
		DatabaseURL:           databaseURL,
		PricingServiceBaseURL: pricingServiceBaseURL,
		SupabaseURL:           supabaseURL,
		SupabaseJWKSURL:       supabaseJWKSURL,
		SupabaseJWTAudience:   supabaseJWTAudience,
		SupabaseJWTIssuer:     supabaseJWTIssuer,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
