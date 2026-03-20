package service

import (
	"context"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
)

type OrderService struct {
	orderRepository    *repository.OrderRepository
	checkoutRepository *repository.CheckoutRepository
}

func NewOrderService(
	orderRepository *repository.OrderRepository,
	checkoutRepository *repository.CheckoutRepository,
) *OrderService {
	return &OrderService{
		orderRepository:    orderRepository,
		checkoutRepository: checkoutRepository,
	}
}

func (s *OrderService) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	return s.orderRepository.Create(ctx, order)
}

func (s *OrderService) CreateWithCheckoutTransaction(ctx context.Context, userID string, order domain.Order) (domain.Order, error) {
	return s.checkoutRepository.Checkout(ctx, userID, order)
}

func (s *OrderService) ListByUserID(ctx context.Context, userID string) ([]domain.Order, error) {
	return s.orderRepository.ListByUserID(ctx, userID)
}

func (s *OrderService) GetByUserIDAndOrderID(ctx context.Context, userID string, orderID string) (domain.Order, bool, error) {
	return s.orderRepository.GetByUserIDAndOrderID(ctx, userID, orderID)
}
