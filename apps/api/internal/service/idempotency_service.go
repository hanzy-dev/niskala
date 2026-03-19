package service

import (
	"context"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
)

type IdempotencyService struct {
	idempotencyRepository *repository.IdempotencyRepository
	orderRepository       *repository.OrderRepository
}

func NewIdempotencyService(
	idempotencyRepository *repository.IdempotencyRepository,
	orderRepository *repository.OrderRepository,
) *IdempotencyService {
	return &IdempotencyService{
		idempotencyRepository: idempotencyRepository,
		orderRepository:       orderRepository,
	}
}

func (s *IdempotencyService) Get(ctx context.Context, userID string, idemKey string) (domain.IdempotencyRecord, bool, error) {
	record, exists, err := s.idempotencyRepository.Get(ctx, userID, idemKey)
	if err != nil || !exists {
		return record, exists, err
	}

	if record.Status == domain.IdempotencyStatusCompleted && record.ResponseOrder != nil && record.ResponseOrder.ID != "" {
		order, ok, err := s.orderRepository.GetByUserIDAndOrderID(ctx, userID, record.ResponseOrder.ID)
		if err != nil {
			return domain.IdempotencyRecord{}, false, err
		}
		if ok {
			record.ResponseOrder = &order
		}
	}

	return record, true, nil
}

func (s *IdempotencyService) StartProcessing(ctx context.Context, userID string, idemKey string) error {
	return s.idempotencyRepository.StartProcessing(ctx, userID, idemKey)
}

func (s *IdempotencyService) Complete(ctx context.Context, userID string, idemKey string, order *domain.Order) error {
	return s.idempotencyRepository.Complete(ctx, userID, idemKey, order.ID)
}

func (s *IdempotencyService) Delete(ctx context.Context, userID string, idemKey string) error {
	return s.idempotencyRepository.Delete(ctx, userID, idemKey)
}
