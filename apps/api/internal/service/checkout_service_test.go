package service

import (
	"context"
	"testing"

	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
)

func TestCheckoutFailsWhenCartIsEmpty(t *testing.T) {
	productRepository := repository.NewProductRepository(nil)
	productService := NewProductService(productRepository)
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
	productRepository := repository.NewProductRepository(nil)
	productService := NewProductService(productRepository)
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
	if err == nil {
		t.Fatalf("expected checkout to fail without seeded products, got order %v", order)
	}
}

func TestCheckoutReplaysCompletedIdempotencyKey(t *testing.T) {
	productRepository := repository.NewProductRepository(nil)
	productService := NewProductService(productRepository)
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

	_, err := checkoutService.Checkout(context.Background(), "user_1", "idem-replay")
	if err == nil {
		t.Fatal("expected first checkout to fail without seeded products, got nil")
	}

	record, exists := idempotencyService.Get("user_1", "idem-replay")
	if !exists {
		t.Fatal("expected idempotency record to exist")
	}

	if record.Status != "processing" {
		t.Fatalf("expected processing status, got %s", record.Status)
	}
}
