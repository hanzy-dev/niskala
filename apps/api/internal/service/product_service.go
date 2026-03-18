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
