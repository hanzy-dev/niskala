package service

import (
	"sync"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
)

type IdempotencyService struct {
	mu      sync.RWMutex
	records map[string]domain.IdempotencyRecord
}

func NewIdempotencyService() *IdempotencyService {
	return &IdempotencyService{
		records: make(map[string]domain.IdempotencyRecord),
	}
}

func (s *IdempotencyService) key(userID string, idemKey string) string {
	return userID + ":" + idemKey
}

func (s *IdempotencyService) Get(userID string, idemKey string) (domain.IdempotencyRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.records[s.key(userID, idemKey)]
	return record, ok
}

func (s *IdempotencyService) StartProcessing(userID string, idemKey string) domain.IdempotencyRecord {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := domain.IdempotencyRecord{
		UserID: userID,
		Key:    idemKey,
		Status: domain.IdempotencyStatusProcessing,
	}
	s.records[s.key(userID, idemKey)] = record
	return record
}

func (s *IdempotencyService) Complete(userID string, idemKey string, order *domain.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[s.key(userID, idemKey)] = domain.IdempotencyRecord{
		UserID:        userID,
		Key:           idemKey,
		Status:        domain.IdempotencyStatusCompleted,
		ResponseOrder: order,
	}
}
