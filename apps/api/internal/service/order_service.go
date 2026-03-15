package service

import (
	"fmt"
	"sync"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
)

type OrderService struct {
	mu      sync.RWMutex
	orders  map[string][]domain.Order
	counter int
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string][]domain.Order),
	}
}

func (s *OrderService) Create(order domain.Order) domain.Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	order.ID = fmt.Sprintf("ord_%d", s.counter)

	s.orders[order.UserID] = append([]domain.Order{order}, s.orders[order.UserID]...)
	return order
}

func (s *OrderService) ListByUserID(userID string) []domain.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.orders[userID]
}

func (s *OrderService) GetByUserIDAndOrderID(userID string, orderID string) (domain.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, order := range s.orders[userID] {
		if order.ID == orderID {
			return order, true
		}
	}

	return domain.Order{}, false
}
