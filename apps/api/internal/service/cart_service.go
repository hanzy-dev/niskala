package service

import (
	"context"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
)

type CartService struct {
	cartRepository *repository.CartRepository
}

func NewCartService(cartRepository *repository.CartRepository) *CartService {
	return &CartService{
		cartRepository: cartRepository,
	}
}

func (s *CartService) GetCart(ctx context.Context, userID string) (domain.Cart, error) {
	return s.cartRepository.GetCart(ctx, userID)
}

func (s *CartService) AddItem(ctx context.Context, userID string, productID string, qty int) (domain.Cart, error) {
	return s.cartRepository.AddItem(ctx, userID, productID, qty)
}

func (s *CartService) UpdateItem(ctx context.Context, userID string, productID string, qty int) (domain.Cart, error) {
	return s.cartRepository.UpdateItem(ctx, userID, productID, qty)
}

func (s *CartService) RemoveItem(ctx context.Context, userID string, productID string) (domain.Cart, error) {
	return s.cartRepository.RemoveItem(ctx, userID, productID)
}

func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	return s.cartRepository.ClearCart(ctx, userID)
}
