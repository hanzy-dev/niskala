package domain

type OrderItem struct {
	ProductID           string `json:"product_id"`
	ProductNameSnapshot string `json:"product_name_snapshot"`
	PriceCents          int64  `json:"price_cents"`
	Qty                 int    `json:"qty"`
}

type Order struct {
	ID                  string      `json:"id"`
	UserID              string      `json:"user_id"`
	Status              string      `json:"status"`
	SubtotalCents       int64       `json:"subtotal_cents"`
	DiscountCents       int64       `json:"discount_cents"`
	TotalCents          int64       `json:"total_cents"`
	PricingFallbackUsed bool        `json:"pricing_fallback_used"`
	Items               []OrderItem `json:"items"`
}
