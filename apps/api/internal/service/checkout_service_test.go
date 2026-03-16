package service

import (
	"context"
	"testing"
)

func TestCheckoutFailsWhenCartIsEmpty(t *testing.T) {
	productService := NewProductService()
	cartService := NewCartService()
	orderService := NewOrderService()
	idempotencyService := NewIdempotencyService()
	pricingService := NewPricingService("http://localhost:8081")

	checkoutService := NewCheckoutService(
		productService,
		cartService,
		orderService,
		idempotencyService,
		pricingService,
	)

	_, err := checkoutService.Checkout(context.Background(), "user_1", "idem-empty")
	if err == nil {
		t.Fatal("expected error for empty cart, got nil")
	}
}

func TestCheckoutCreatesOrderAndClearsCart(t *testing.T) {
	productService := NewProductService()
	cartService := NewCartService()
	orderService := NewOrderService()
	idempotencyService := NewIdempotencyService()
	pricingService := NewPricingService("http://localhost:8081")

	cartService.AddItem("user_1", "prod_1", 2)

	checkoutService := NewCheckoutService(
		productService,
		cartService,
		orderService,
		idempotencyService,
		pricingService,
	)

	order, err := checkoutService.Checkout(context.Background(), "user_1", "idem-ok")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if order.ID == "" {
		t.Fatal("expected order id to be set")
	}

	cart := cartService.GetCart("user_1")
	if len(cart.Items) != 0 {
		t.Fatal("expected cart to be cleared after checkout")
	}
}

func TestCheckoutReplaysCompletedIdempotencyKey(t *testing.T) {
	productService := NewProductService()
	cartService := NewCartService()
	orderService := NewOrderService()
	idempotencyService := NewIdempotencyService()
	pricingService := NewPricingService("http://localhost:8081")

	cartService.AddItem("user_1", "prod_1", 1)

	checkoutService := NewCheckoutService(
		productService,
		cartService,
		orderService,
		idempotencyService,
		pricingService,
	)

	firstOrder, err := checkoutService.Checkout(context.Background(), "user_1", "idem-replay")
	if err != nil {
		t.Fatalf("expected first checkout to succeed, got %v", err)
	}

	secondOrder, err := checkoutService.Checkout(context.Background(), "user_1", "idem-replay")
	if err != nil {
		t.Fatalf("expected replay checkout to succeed, got %v", err)
	}

	if firstOrder.ID != secondOrder.ID {
		t.Fatalf("expected replayed order id %s, got %s", firstOrder.ID, secondOrder.ID)
	}
}
