package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hanzy-dev/niskala/apps/api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IdempotencyRepository struct {
	db *pgxpool.Pool
}

func NewIdempotencyRepository(db *pgxpool.Pool) *IdempotencyRepository {
	return &IdempotencyRepository{
		db: db,
	}
}

func (r *IdempotencyRepository) Get(ctx context.Context, userID string, key string) (domain.IdempotencyRecord, bool, error) {
	if r.db == nil {
		return domain.IdempotencyRecord{}, false, fmt.Errorf("database connection is not available")
	}

	var record domain.IdempotencyRecord
	var orderID *string

	err := r.db.QueryRow(ctx, `
		SELECT user_id, key, status, order_id
		FROM idempotency_keys
		WHERE user_id = $1 AND key = $2
		LIMIT 1
	`, userID, key).Scan(
		&record.UserID,
		&record.Key,
		&record.Status,
		&orderID,
	)
	if err != nil {
		return domain.IdempotencyRecord{}, false, nil
	}

	if orderID != nil {
		record.ResponseOrder = &domain.Order{ID: *orderID}
	}

	return record, true, nil
}

func (r *IdempotencyRepository) StartProcessing(ctx context.Context, userID string, key string) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not available")
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO idempotency_keys (
			user_id,
			key,
			status,
			order_id,
			response_json,
			created_at,
			expires_at
		)
		VALUES ($1, $2, 'processing', NULL, NULL, NOW(), $3)
		ON CONFLICT (user_id, key) DO NOTHING
	`, userID, key, time.Now().Add(24*time.Hour))
	if err != nil {
		return fmt.Errorf("insert idempotency processing record: %w", err)
	}

	return nil
}

func (r *IdempotencyRepository) Complete(ctx context.Context, userID string, key string, orderID string) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not available")
	}

	_, err := r.db.Exec(ctx, `
		UPDATE idempotency_keys
		SET
			status = 'completed',
			order_id = $3
		WHERE user_id = $1 AND key = $2
	`, userID, key, orderID)
	if err != nil {
		return fmt.Errorf("complete idempotency record: %w", err)
	}

	return nil
}

func (r *IdempotencyRepository) Delete(ctx context.Context, userID string, key string) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not available")
	}

	_, err := r.db.Exec(ctx, `
		DELETE FROM idempotency_keys
		WHERE user_id = $1 AND key = $2
	`, userID, key)
	if err != nil {
		return fmt.Errorf("delete idempotency record: %w", err)
	}

	return nil
}
