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

func (r *ProductRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	if r.db == nil {
		return domain.Product{}, fmt.Errorf("database connection is not available")
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO products (
			id, name, description, price_cents, stock, category, image_url, is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`,
		product.ID,
		product.Name,
		product.Description,
		product.PriceCents,
		product.Stock,
		product.Category,
		product.ImageURL,
		product.IsActive,
	)
	if err != nil {
		return domain.Product{}, fmt.Errorf("create product: %w", err)
	}

	return product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	if r.db == nil {
		return domain.Product{}, fmt.Errorf("database connection is not available")
	}

	_, err := r.db.Exec(ctx, `
		UPDATE products
		SET
			name = $2,
			description = $3,
			price_cents = $4,
			category = $5,
			image_url = $6,
			is_active = $7,
			updated_at = NOW()
		WHERE id = $1
	`,
		product.ID,
		product.Name,
		product.Description,
		product.PriceCents,
		product.Category,
		product.ImageURL,
		product.IsActive,
	)
	if err != nil {
		return domain.Product{}, fmt.Errorf("update product: %w", err)
	}

	return product, nil
}

func (r *ProductRepository) UpdateStock(ctx context.Context, productID string, stock int) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not available")
	}

	_, err := r.db.Exec(ctx, `
		UPDATE products
		SET stock = $2, updated_at = NOW()
		WHERE id = $1
	`, productID, stock)
	if err != nil {
		return fmt.Errorf("update product stock: %w", err)
	}

	return nil
}

func (r *ProductRepository) DecrementStock(ctx context.Context, productID string, qty int) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not available")
	}

	commandTag, err := r.db.Exec(ctx, `
		UPDATE products
		SET stock = stock - $2, updated_at = NOW()
		WHERE id = $1 AND stock >= $2
	`, productID, qty)
	if err != nil {
		return fmt.Errorf("decrement product stock: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("insufficient stock")
	}

	return nil
}
