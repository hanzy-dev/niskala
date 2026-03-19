package service

import (
	"context"
	"testing"

	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
)

func TestCheckoutFailsWhenCartIsEmpty(t *testing.T) {
	productRepository := repository.NewProductRepository(nil)
	cartRepository := repository.NewCartRepository(nil)
	orderRepository := repository.NewOrderRepository(nil)
	idempotencyRepository := repository.NewIdempotencyRepository(nil)

	productService := NewProductService(productRepository)
	cartService := NewCartService(cartRepository)
	orderService := NewOrderService(orderRepository)
	idempotencyService := NewIdempotencyService(idempotencyRepository, orderRepository)
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

func TestCheckoutFailsWithoutSeededProducts(t *testing.T) {
	productRepository := repository.NewProductRepository(nil)
	cartRepository := repository.NewCartRepository(nil)
	orderRepository := repository.NewOrderRepository(nil)
	idempotencyRepository := repository.NewIdempotencyRepository(nil)

	productService := NewProductService(productRepository)
	cartService := NewCartService(cartRepository)
	orderService := NewOrderService(orderRepository)
	idempotencyService := NewIdempotencyService(idempotencyRepository, orderRepository)
	pricingService := NewPricingService("http://localhost:8081")

	checkoutService := NewCheckoutService(
		productService,
		cartService,
		orderService,
		idempotencyService,
		pricingService,
	)

	_, err := checkoutService.Checkout(context.Background(), "user_1", "idem-ok")
	if err == nil {
		t.Fatal("expected checkout to fail without seeded products, got nil")
	}
}
