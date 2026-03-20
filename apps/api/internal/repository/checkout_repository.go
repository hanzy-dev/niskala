package repository

import (
	"context"
	"fmt"
	"sort"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CheckoutRepository struct {
	db *pgxpool.Pool
}

func NewCheckoutRepository(db *pgxpool.Pool) *CheckoutRepository {
	return &CheckoutRepository{
		db: db,
	}
}

func (r *CheckoutRepository) Checkout(ctx context.Context, userID string, order domain.Order) (domain.Order, error) {
	if r.db == nil {
		return domain.Order{}, fmt.Errorf("database connection is not available")
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Order{}, fmt.Errorf("begin checkout transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	sortedItems := make([]domain.OrderItem, len(order.Items))
	copy(sortedItems, order.Items)
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].ProductID < sortedItems[j].ProductID
	})

	for _, item := range sortedItems {
		var stock int
		err := tx.QueryRow(ctx, `
			SELECT stock
			FROM products
			WHERE id = $1
			FOR UPDATE
		`, item.ProductID).Scan(&stock)
		if err != nil {
			return domain.Order{}, fmt.Errorf("lock product %s: %w", item.ProductID, err)
		}

		if stock < item.Qty {
			return domain.Order{}, fmt.Errorf("insufficient stock")
		}
	}

	for _, item := range sortedItems {
		_, err := tx.Exec(ctx, `
			UPDATE products
			SET stock = stock - $2, updated_at = NOW()
			WHERE id = $1
		`, item.ProductID, item.Qty)
		if err != nil {
			return domain.Order{}, fmt.Errorf("decrement stock for %s: %w", item.ProductID, err)
		}
	}

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

	for _, item := range sortedItems {
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

	var cartID string
	err = tx.QueryRow(ctx, `
		SELECT id
		FROM carts
		WHERE user_id = $1
		LIMIT 1
	`, userID).Scan(&cartID)
	if err == nil {
		_, err = tx.Exec(ctx, `
			DELETE FROM cart_items
			WHERE cart_id = $1
		`, cartID)
		if err != nil {
			return domain.Order{}, fmt.Errorf("clear cart: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Order{}, fmt.Errorf("commit checkout transaction: %w", err)
	}

	order.ID = orderID
	order.Items = sortedItems
	return order, nil
}
