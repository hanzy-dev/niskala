package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
)

type PricingService struct {
	baseURL    string
	httpClient *http.Client
}

type pricingQuoteRequest struct {
	Items      []pricingQuoteItem `json:"items"`
	CouponCode *string            `json:"coupon_code"`
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
	requestItems := make([]pricingQuoteItem, 0, len(items))
	for _, item := range items {
		requestItems = append(requestItems, pricingQuoteItem{
			ProductID:  item.ProductID,
			PriceCents: item.PriceCents,
			Qty:        item.Qty,
		})
	}

	payload := pricingQuoteRequest{
		Items:      requestItems,
		CouponCode: nil,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, 0, 0, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL+"/pricing/quote",
		bytes.NewReader(body),
	)
	if err != nil {
		return 0, 0, 0, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return 0, 0, 0, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 {
		return 0, 0, 0, err
	}

	var quote pricingQuoteResponse
	if err := json.NewDecoder(response.Body).Decode(&quote); err != nil {
		return 0, 0, 0, err
	}

	return quote.SubtotalCents, quote.DiscountCents, quote.TotalCents, nil
}
