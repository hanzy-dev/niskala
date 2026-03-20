package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/redis/go-redis/v9"
)

type ProductCacheService struct {
	rdb *redis.Client
}

func NewProductCacheService(rdb *redis.Client) *ProductCacheService {
	return &ProductCacheService{
		rdb: rdb,
	}
}

func (s *ProductCacheService) listKey() string {
	return "products:list:active"
}

func (s *ProductCacheService) productKey(productID string) string {
	return fmt.Sprintf("products:detail:%s", productID)
}

func (s *ProductCacheService) GetProductList(ctx context.Context) ([]domain.Product, bool) {
	if s.rdb == nil {
		return nil, false
	}

	value, err := s.rdb.Get(ctx, s.listKey()).Result()
	if err != nil {
		return nil, false
	}

	var products []domain.Product
	if err := json.Unmarshal([]byte(value), &products); err != nil {
		return nil, false
	}

	return products, true
}

func (s *ProductCacheService) SetProductList(ctx context.Context, products []domain.Product) {
	if s.rdb == nil {
		return
	}

	body, err := json.Marshal(products)
	if err != nil {
		return
	}

	_ = s.rdb.Set(ctx, s.listKey(), body, 60*time.Second).Err()
}

func (s *ProductCacheService) GetProduct(ctx context.Context, productID string) (domain.Product, bool) {
	if s.rdb == nil {
		return domain.Product{}, false
	}

	value, err := s.rdb.Get(ctx, s.productKey(productID)).Result()
	if err != nil {
		return domain.Product{}, false
	}

	var product domain.Product
	if err := json.Unmarshal([]byte(value), &product); err != nil {
		return domain.Product{}, false
	}

	return product, true
}

func (s *ProductCacheService) SetProduct(ctx context.Context, product domain.Product) {
	if s.rdb == nil {
		return
	}

	body, err := json.Marshal(product)
	if err != nil {
		return
	}

	_ = s.rdb.Set(ctx, s.productKey(product.ID), body, 60*time.Second).Err()
}

func (s *ProductCacheService) InvalidateProduct(ctx context.Context, productID string) {
	if s.rdb == nil {
		return
	}

	_ = s.rdb.Del(ctx, s.listKey(), s.productKey(productID)).Err()
}

func (s *ProductCacheService) InvalidateProductList(ctx context.Context) {
	if s.rdb == nil {
		return
	}

	_ = s.rdb.Del(ctx, s.listKey()).Err()
}
