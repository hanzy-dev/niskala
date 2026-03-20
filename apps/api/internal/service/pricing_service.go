package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
)

type PricingService struct {
	baseURL    string
	httpClient *http.Client
}

type pricingQuoteRequest struct {
	Items []pricingQuoteItem `json:"items"`
}

type pricingQuoteItem struct {
	ProductID  string `json:"product_id"`
	PriceCents int64  `json:"price_cents"`
	Qty        int    `json:"qty"`
}

type pricingQuoteResponse struct {
	SubtotalCents int64 `json:"subtotal_cents"`
	DiscountCents int64 `json:"discount_cents"`
	TotalCents    int64 `json:"total_cents"`
}

func NewPricingService(baseURL string) *PricingService {
	return &PricingService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 250 * time.Millisecond,
		},
	}
}

func (s *PricingService) Quote(ctx context.Context, items []domain.OrderItem) (int64, int64, int64, error) {
	if s.baseURL == "" {
		return 0, 0, 0, fmt.Errorf("pricing service base url is empty")
	}

	requestItems := make([]pricingQuoteItem, 0, len(items))
	for _, item := range items {
		requestItems = append(requestItems, pricingQuoteItem{
			ProductID:  item.ProductID,
			PriceCents: item.PriceCents,
			Qty:        item.Qty,
		})
	}

	payload := pricingQuoteRequest{
		Items: requestItems,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("marshal pricing request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL+"/pricing/quote",
		bytes.NewReader(body),
	)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("build pricing request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("call pricing service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, 0, 0, fmt.Errorf("pricing service returned status %d", resp.StatusCode)
	}

	var quote pricingQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return 0, 0, 0, fmt.Errorf("decode pricing response: %w", err)
	}

	return quote.SubtotalCents, quote.DiscountCents, quote.TotalCents, nil
}
