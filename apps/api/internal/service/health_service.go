package service

import (
	"context"
	"net/http"
	"time"

	"github.com/hanzy-dev/niskala/apps/api/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthService struct {
	db                    *pgxpool.Pool
	pricingServiceBaseURL string
	httpClient            *http.Client
}

type HealthStatus struct {
	Status            string `json:"status"`
	Service           string `json:"service"`
	Database          string `json:"database"`
	PricingService    string `json:"pricing_service"`
	CheckoutAvailable bool   `json:"checkout_available"`
}

func NewHealthService(db *pgxpool.Pool, pricingServiceBaseURL string) *HealthService {
	return &HealthService{
		db:                    db,
		pricingServiceBaseURL: pricingServiceBaseURL,
		httpClient: &http.Client{
			Timeout: 200 * time.Millisecond,
		},
	}
}

func (s *HealthService) GetStatus(ctx context.Context) HealthStatus {
	databaseStatus := "down"
	pricingStatus := "down"

	if s.db != nil && database.Ping(ctx, s.db) == nil {
		databaseStatus = "ok"
	}

	response, err := s.httpClient.Get(s.pricingServiceBaseURL + "/health")
	if err == nil && response != nil {
		_ = response.Body.Close()
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			pricingStatus = "ok"
		}
	}

	status := "ok"
	checkoutAvailable := true

	if databaseStatus != "ok" {
		status = "down"
		checkoutAvailable = false
	} else if pricingStatus != "ok" {
		status = "degraded"
	}

	return HealthStatus{
		Status:            status,
		Service:           "api",
		Database:          databaseStatus,
		PricingService:    pricingStatus,
		CheckoutAvailable: checkoutAvailable,
	}
}
