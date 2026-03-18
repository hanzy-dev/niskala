package repository

import (
	"context"
	"fmt"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) ensureCart(ctx context.Context, userID string) (string, error) {
	if r.db == nil {
		return "", fmt.Errorf("database connection is not available")
	}

	var cartID string
	err := r.db.QueryRow(ctx, `
		INSERT INTO carts (id, user_id)
		VALUES ('cart_' || replace(gen_random_uuid()::text, '-', ''), $1)
		ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
		RETURNING id
	`, userID).Scan(&cartID)
	if err != nil {
		return "", fmt.Errorf("ensure cart: %w", err)
	}

	return cartID, nil
}

func (r *CartRepository) GetCart(ctx context.Context, userID string) (domain.Cart, error) {
	if r.db == nil {
		return domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}, nil
	}

	cartID, err := r.ensureCart(ctx, userID)
	if err != nil {
		return domain.Cart{}, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT product_id, qty
		FROM cart_items
		WHERE cart_id = $1
		ORDER BY product_id ASC
	`, cartID)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("query cart items: %w", err)
	}
	defer rows.Close()

	items := make([]domain.CartItem, 0)
	for rows.Next() {
		var item domain.CartItem
		if err := rows.Scan(&item.ProductID, &item.Qty); err != nil {
			return domain.Cart{}, fmt.Errorf("scan cart item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return domain.Cart{}, fmt.Errorf("iterate cart items: %w", err)
	}

	return domain.Cart{
		UserID: userID,
		Items:  items,
	}, nil
}

func (r *CartRepository) AddItem(ctx context.Context, userID string, productID string, qty int) (domain.Cart, error) {
	if r.db == nil {
		return domain.Cart{}, fmt.Errorf("database connection is not available")
	}

	cartID, err := r.ensureCart(ctx, userID)
	if err != nil {
		return domain.Cart{}, err
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO cart_items (cart_id, product_id, qty)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET qty = cart_items.qty + EXCLUDED.qty, updated_at = NOW()
	`, cartID, productID, qty)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("add cart item: %w", err)
	}

	return r.GetCart(ctx, userID)
}

func (r *CartRepository) UpdateItem(ctx context.Context, userID string, productID string, qty int) (domain.Cart, error) {
	if r.db == nil {
		return domain.Cart{}, fmt.Errorf("database connection is not available")
	}

	cartID, err := r.ensureCart(ctx, userID)
	if err != nil {
		return domain.Cart{}, err
	}

	if qty <= 0 {
		_, err = r.db.Exec(ctx, `
			DELETE FROM cart_items
			WHERE cart_id = $1 AND product_id = $2
		`, cartID, productID)
		if err != nil {
			return domain.Cart{}, fmt.Errorf("delete cart item during update: %w", err)
		}

		return r.GetCart(ctx, userID)
	}

	_, err = r.db.Exec(ctx, `
		UPDATE cart_items
		SET qty = $3, updated_at = NOW()
		WHERE cart_id = $1 AND product_id = $2
	`, cartID, productID, qty)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("update cart item: %w", err)
	}

	return r.GetCart(ctx, userID)
}

func (r *CartRepository) RemoveItem(ctx context.Context, userID string, productID string) (domain.Cart, error) {
	if r.db == nil {
		return domain.Cart{}, fmt.Errorf("database connection is not available")
	}

	cartID, err := r.ensureCart(ctx, userID)
	if err != nil {
		return domain.Cart{}, err
	}

	_, err = r.db.Exec(ctx, `
		DELETE FROM cart_items
		WHERE cart_id = $1 AND product_id = $2
	`, cartID, productID)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("remove cart item: %w", err)
	}

	return r.GetCart(ctx, userID)
}

func (r *CartRepository) ClearCart(ctx context.Context, userID string) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not available")
	}

	cartID, err := r.ensureCart(ctx, userID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		DELETE FROM cart_items
		WHERE cart_id = $1
	`, cartID)
	if err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}

	return nil
}
