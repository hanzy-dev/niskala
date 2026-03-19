package service

import (
	"context"
	"errors"
	"strings"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
)

var (
	ErrEmptyCart             = errors.New("empty cart")
	ErrProductNotFound       = errors.New("product not found")
	ErrInsufficientStock     = errors.New("insufficient stock")
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

	cart, err := s.cartService.GetCart(ctx, userID)
	if err != nil {
		return domain.Order{}, err
	}

	if len(cart.Items) == 0 {
		return domain.Order{}, ErrEmptyCart
	}

	var subtotalCents int64
	orderItems := make([]domain.OrderItem, 0, len(cart.Items))

	for _, cartItem := range cart.Items {
		product, ok, err := s.productService.GetByID(ctx, cartItem.ProductID)
		if err != nil {
			return domain.Order{}, err
		}
		if !ok {
			return domain.Order{}, ErrProductNotFound
		}

		if product.Stock < cartItem.Qty {
			return domain.Order{}, ErrInsufficientStock
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

	for _, item := range orderItems {
		if err := s.productService.DecrementStock(ctx, item.ProductID, item.Qty); err != nil {
			if strings.Contains(err.Error(), "insufficient stock") {
				return domain.Order{}, ErrInsufficientStock
			}
			return domain.Order{}, err
		}
	}

	order, err := s.orderService.Create(ctx, domain.Order{
		UserID:              userID,
		Status:              "created",
		SubtotalCents:       subtotalCents,
		DiscountCents:       discountCents,
		TotalCents:          totalCents,
		PricingFallbackUsed: pricingFallbackUsed,
		Items:               orderItems,
	})
	if err != nil {
		return domain.Order{}, err
	}

	if err := s.cartService.ClearCart(ctx, userID); err != nil {
		return domain.Order{}, err
	}

	s.idempotencyService.Complete(userID, idemKey, &order)

	return order, nil
}
