package service

import (
	"context"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
)

type ProductService struct {
	productRepository *repository.ProductRepository
}

func NewProductService(productRepository *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) List(ctx context.Context) ([]domain.Product, error) {
	return s.productRepository.ListActive(ctx)
}

func (s *ProductService) GetByID(ctx context.Context, id string) (domain.Product, bool, error) {
	return s.productRepository.GetByID(ctx, id)
}

func (s *ProductService) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	return s.productRepository.Create(ctx, product)
}

func (s *ProductService) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	return s.productRepository.Update(ctx, product)
}

func (s *ProductService) UpdateStock(ctx context.Context, productID string, stock int) error {
	return s.productRepository.UpdateStock(ctx, productID, stock)
}

func (s *ProductService) DecrementStock(ctx context.Context, productID string, qty int) error {
	return s.productRepository.DecrementStock(ctx, productID, qty)
}
