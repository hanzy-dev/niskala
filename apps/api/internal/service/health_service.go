package service

import (
	"net/http"
	"time"
)

type HealthService struct {
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

func NewHealthService(pricingServiceBaseURL string) *HealthService {
	return &HealthService{
		pricingServiceBaseURL: pricingServiceBaseURL,
		httpClient: &http.Client{
			Timeout: 200 * time.Millisecond,
		},
	}
}

func (s *HealthService) GetStatus() HealthStatus {
	pricingStatus := "down"

	response, err := s.httpClient.Get(s.pricingServiceBaseURL + "/health")
	if err == nil && response != nil {
		_ = response.Body.Close()
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			pricingStatus = "ok"
		}
	}

	status := "ok"
	checkoutAvailable := true

	if pricingStatus != "ok" {
		status = "degraded"
	}

	return HealthStatus{
		Status:            status,
		Service:           "api",
		Database:          "unknown",
		PricingService:    pricingStatus,
		CheckoutAvailable: checkoutAvailable,
	}
}
