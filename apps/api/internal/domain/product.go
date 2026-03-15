package domain

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int64  `json:"price_cents"`
	Stock       int    `json:"stock"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}
