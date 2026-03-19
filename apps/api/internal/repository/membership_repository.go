package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MembershipRepository struct {
	db *pgxpool.Pool
}

func NewMembershipRepository(db *pgxpool.Pool) *MembershipRepository {
	return &MembershipRepository{
		db: db,
	}
}

func (r *MembershipRepository) GetRoleByUserID(ctx context.Context, userID string) (string, bool, error) {
	if r.db == nil {
		return "", false, fmt.Errorf("database connection is not available")
	}

	var role string
	err := r.db.QueryRow(ctx, `
		SELECT role
		FROM memberships
		WHERE user_id = $1
		LIMIT 1
	`, userID).Scan(&role)
	if err != nil {
		return "", false, nil
	}

	return role, true, nil
}
