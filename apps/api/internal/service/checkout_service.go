package service

import (
	"context"
	"errors"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
)

var (
	ErrEmptyCart             = errors.New("empty cart")
	ErrProductNotFound       = errors.New("product not found")
	ErrIdempotencyInProgress = errors.New("idempotency in progress")
	ErrMissingIdempotencyKey = errors.New("missing idempotency key")
)

type CheckoutService struct {
	productService     *ProductService
	cartService        *CartService
	orderService       *OrderService
	idempotencyService *IdempotencyService
	pricingService     *PricingService
}

func NewCheckoutService(
	productService *ProductService,
	cartService *CartService,
	orderService *OrderService,
	idempotencyService *IdempotencyService,
	pricingService *PricingService,
) *CheckoutService {
	return &CheckoutService{
		productService:     productService,
		cartService:        cartService,
		orderService:       orderService,
		idempotencyService: idempotencyService,
		pricingService:     pricingService,
	}
}

func (s *CheckoutService) Checkout(ctx context.Context, userID string, idemKey string) (domain.Order, error) {
	if idemKey == "" {
		return domain.Order{}, ErrMissingIdempotencyKey
	}

	if record, exists := s.idempotencyService.Get(userID, idemKey); exists {
		if record.Status == domain.IdempotencyStatusCompleted && record.ResponseOrder != nil {
			return *record.ResponseOrder, nil
		}

		if record.Status == domain.IdempotencyStatusProcessing {
			return domain.Order{}, ErrIdempotencyInProgress
		}
	}

	s.idempotencyService.StartProcessing(userID, idemKey)

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

	discountCents := int64(0)
	totalCents := subtotalCents
	pricingFallbackUsed := false

	quotedSubtotal, quotedDiscount, quotedTotal, err := s.pricingService.Quote(ctx, orderItems)
	if err == nil {
		subtotalCents = quotedSubtotal
		discountCents = quotedDiscount
		totalCents = quotedTotal
	} else {
		pricingFallbackUsed = true
	}

	order := s.orderService.Create(domain.Order{
		UserID:              userID,
		Status:              "created",
		SubtotalCents:       subtotalCents,
		DiscountCents:       discountCents,
		TotalCents:          totalCents,
		PricingFallbackUsed: pricingFallbackUsed,
		Items:               orderItems,
	})

	s.cartService.ClearCart(userID)
	s.idempotencyService.Complete(userID, idemKey, &order)

	return order, nil
}
