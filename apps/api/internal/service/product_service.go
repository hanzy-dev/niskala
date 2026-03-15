package service

import "github.com/hanzy-dev/niskala/apps/api/internal/domain"

type ProductService struct {
	products []domain.Product
}

func NewProductService() *ProductService {
	return &ProductService{
		products: []domain.Product{
			{
				ID:          "prod_1",
				Name:        "Niskala Journal",
				Description: "A minimal notebook designed for calm work.",
				PriceCents:  150000,
				Stock:       12,
				Category:    "stationery",
				ImageURL:    "",
				IsActive:    true,
			},
			{
				ID:          "prod_2",
				Name:        "Niskala Bottle",
				Description: "A stainless bottle with a clean silhouette.",
				PriceCents:  220000,
				Stock:       8,
				Category:    "lifestyle",
				ImageURL:    "",
				IsActive:    true,
			},
		},
	}
}

func (s *ProductService) List() []domain.Product {
	return s.products
}

func (s *ProductService) GetByID(id string) (domain.Product, bool) {
	for _, product := range s.products {
		if product.ID == id {
			return product, true
		}
	}

	return domain.Product{}, false
}
