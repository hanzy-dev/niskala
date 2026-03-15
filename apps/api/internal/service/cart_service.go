package service

import (
	"sync"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
)

type CartService struct {
	mu    sync.RWMutex
	carts map[string]domain.Cart
}

func NewCartService() *CartService {
	return &CartService{
		carts: make(map[string]domain.Cart),
	}
}

func (s *CartService) GetCart(userID string) domain.Cart {
	s.mu.RLock()
	cart, exists := s.carts[userID]
	s.mu.RUnlock()

	if !exists {
		return domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}
	}

	return cart
}

func (s *CartService) AddItem(userID string, productID string, qty int) domain.Cart {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, exists := s.carts[userID]
	if !exists {
		cart = domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}
	}

	for index, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[index].Qty += qty
			s.carts[userID] = cart
			return cart
		}
	}

	cart.Items = append(cart.Items, domain.CartItem{
		ProductID: productID,
		Qty:       qty,
	})

	s.carts[userID] = cart
	return cart
}

func (s *CartService) UpdateItem(userID string, productID string, qty int) domain.Cart {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, exists := s.carts[userID]
	if !exists {
		return domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}
	}

	for index, item := range cart.Items {
		if item.ProductID == productID {
			if qty <= 0 {
				cart.Items = append(cart.Items[:index], cart.Items[index+1:]...)
			} else {
				cart.Items[index].Qty = qty
			}

			s.carts[userID] = cart
			return cart
		}
	}

	return cart
}

func (s *CartService) RemoveItem(userID string, productID string) domain.Cart {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, exists := s.carts[userID]
	if !exists {
		return domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}
	}

	for index, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:index], cart.Items[index+1:]...)
			s.carts[userID] = cart
			return cart
		}
	}

	return cart
}

func (s *CartService) ClearCart(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.carts[userID] = domain.Cart{
		UserID: userID,
		Items:  []domain.CartItem{},
	}
}
