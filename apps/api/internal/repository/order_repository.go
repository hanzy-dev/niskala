package repository

import (
	"context"
	"fmt"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	if r.db == nil {
		return domain.Order{}, fmt.Errorf("database connection is not available")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Order{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var orderID string
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (
			id,
			user_id,
			status,
			subtotal_cents,
			discount_cents,
			total_cents,
			pricing_fallback_used,
			metadata_json,
			created_at
		)
		VALUES (
			'ord_' || replace(gen_random_uuid()::text, '-', ''),
			$1, $2, $3, $4, $5, $6, '{}'::jsonb, NOW()
		)
		RETURNING id
	`,
		order.UserID,
		order.Status,
		order.SubtotalCents,
		order.DiscountCents,
		order.TotalCents,
		order.PricingFallbackUsed,
	).Scan(&orderID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("insert order: %w", err)
	}

	for _, item := range order.Items {
		_, err := tx.Exec(ctx, `
			INSERT INTO order_items (
				order_id,
				product_id,
				product_name_snapshot,
				price_cents,
				qty
			)
			VALUES ($1, $2, $3, $4, $5)
		`,
			orderID,
			item.ProductID,
			item.ProductNameSnapshot,
			item.PriceCents,
			item.Qty,
		)
		if err != nil {
			return domain.Order{}, fmt.Errorf("insert order item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Order{}, fmt.Errorf("commit order transaction: %w", err)
	}

	order.ID = orderID
	return order, nil
}

func (r *OrderRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Order, error) {
	if r.db == nil {
		return []domain.Order{}, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, status, subtotal_cents, discount_cents, total_cents, pricing_fallback_used
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.SubtotalCents,
			&order.DiscountCents,
			&order.TotalCents,
			&order.PricingFallbackUsed,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		order.Items = []domain.OrderItem{}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) GetByUserIDAndOrderID(ctx context.Context, userID string, orderID string) (domain.Order, bool, error) {
	if r.db == nil {
		return domain.Order{}, false, nil
	}

	var order domain.Order
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, status, subtotal_cents, discount_cents, total_cents, pricing_fallback_used
		FROM orders
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`, orderID, userID).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.SubtotalCents,
		&order.DiscountCents,
		&order.TotalCents,
		&order.PricingFallbackUsed,
	)
	if err != nil {
		return domain.Order{}, false, nil
	}

	itemRows, err := r.db.Query(ctx, `
		SELECT product_id, product_name_snapshot, price_cents, qty
		FROM order_items
		WHERE order_id = $1
		ORDER BY product_id ASC
	`, orderID)
	if err != nil {
		return domain.Order{}, false, fmt.Errorf("query order items: %w", err)
	}
	defer itemRows.Close()

	items := make([]domain.OrderItem, 0)
	for itemRows.Next() {
		var item domain.OrderItem
		if err := itemRows.Scan(
			&item.ProductID,
			&item.ProductNameSnapshot,
			&item.PriceCents,
			&item.Qty,
		); err != nil {
			return domain.Order{}, false, fmt.Errorf("scan order item: %w", err)
		}

		items = append(items, item)
	}

	if err := itemRows.Err(); err != nil {
		return domain.Order{}, false, fmt.Errorf("iterate order items: %w", err)
	}

	order.Items = items
	return order, true, nil
}
