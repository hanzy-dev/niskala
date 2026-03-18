package repository

import (
	"context"
	"fmt"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) ListActive(ctx context.Context) ([]domain.Product, error) {
	if r.db == nil {
		return []domain.Product{}, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, name, description, price_cents, stock, category, image_url, is_active
		FROM products
		WHERE is_active = TRUE
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query active products: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.PriceCents,
			&product.Stock,
			&product.Category,
			&product.ImageURL,
			&product.IsActive,
		); err != nil {
			return nil, fmt.Errorf("scan active product: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active products: %w", err)
	}

	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (domain.Product, bool, error) {
	if r.db == nil {
		return domain.Product{}, false, nil
	}

	var product domain.Product
	err := r.db.QueryRow(ctx, `
		SELECT id, name, description, price_cents, stock, category, image_url, is_active
		FROM products
		WHERE id = $1
		LIMIT 1
	`, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.PriceCents,
		&product.Stock,
		&product.Category,
		&product.ImageURL,
		&product.IsActive,
	)
	if err != nil {
		return domain.Product{}, false, nil
	}

	return product, true, nil
}
