package service

import (
	"errors"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
)

var (
	ErrEmptyCart       = errors.New("empty cart")
	ErrProductNotFound = errors.New("product not found")
)

type CheckoutService struct {
	productService *ProductService
	cartService    *CartService
	orderService   *OrderService
}

func NewCheckoutService(
	productService *ProductService,
	cartService *CartService,
	orderService *OrderService,
) *CheckoutService {
	return &CheckoutService{
		productService: productService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

func (s *CheckoutService) Checkout(userID string) (domain.Order, error) {
	cart := s.cartService.GetCart(userID)
	if len(cart.Items) == 0 {
		return domain.Order{}, ErrEmptyCart
	}

	var subtotalCents int64
	orderItems := make([]domain.OrderItem, 0, len(cart.Items))

	for _, cartItem := range cart.Items {
		product, ok := s.productService.GetByID(cartItem.ProductID)
		if !ok {
			return domain.Order{}, ErrProductNotFound
		}

		subtotalCents += product.PriceCents * int64(cartItem.Qty)

		orderItems = append(orderItems, domain.OrderItem{
			ProductID:           product.ID,
			ProductNameSnapshot: product.Name,
			PriceCents:          product.PriceCents,
			Qty:                 cartItem.Qty,
		})
	}

	order := s.orderService.Create(domain.Order{
		UserID:              userID,
		Status:              "created",
		SubtotalCents:       subtotalCents,
		DiscountCents:       0,
		TotalCents:          subtotalCents,
		PricingFallbackUsed: false,
		Items:               orderItems,
	})

	s.cartService.ClearCart(userID)

	return order, nil
}
