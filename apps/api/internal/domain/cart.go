package domain

type CartItem struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}
