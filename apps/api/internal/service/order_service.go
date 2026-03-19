package service

import (
	"context"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
)

type OrderService struct {
	orderRepository *repository.OrderRepository
}

func NewOrderService(orderRepository *repository.OrderRepository) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
	}
}

func (s *OrderService) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	return s.orderRepository.Create(ctx, order)
}

func (s *OrderService) ListByUserID(ctx context.Context, userID string) ([]domain.Order, error) {
	return s.orderRepository.ListByUserID(ctx, userID)
}

func (s *OrderService) GetByUserIDAndOrderID(ctx context.Context, userID string, orderID string) (domain.Order, bool, error) {
	return s.orderRepository.GetByUserIDAndOrderID(ctx, userID, orderID)
}
