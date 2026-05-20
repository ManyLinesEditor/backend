package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{pool: pool}
}

// Upsert creates or renews a subscription.
func (r *SubscriptionRepo) Upsert(ctx context.Context, userID, featureID string, until time.Time) error {
	const q = `
		INSERT INTO subscriptions (user_id, feature_id, until)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, feature_id)
		DO UPDATE SET until = GREATEST(subscriptions.until, EXCLUDED.until)`

	_, err := r.pool.Exec(ctx, q, userID, featureID, until)
	if err != nil {
		return fmt.Errorf("subscription upsert: %w", err)
	}
	return nil
}

// ListActive returns the user's active subscriptions.
func (r *SubscriptionRepo) ListActive(ctx context.Context, userID string) ([]*models.Subscription, error) {
	const q = `
		SELECT user_id, feature_id, until
		FROM subscriptions
		WHERE user_id = $1 AND until > now()`

	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Subscription
	for rows.Next() {
		s := &models.Subscription{}
		if err := rows.Scan(&s.UserID, &s.FeatureID, &s.Until); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// HasActive checks for a specific feature.
func (r *SubscriptionRepo) HasActive(ctx context.Context, userID, featureID string) (bool, error) {
	const q = `SELECT 1 FROM subscriptions WHERE user_id=$1 AND feature_id=$2 AND until > now()`
	var dummy int
	err := r.pool.QueryRow(ctx, q, userID, featureID).Scan(&dummy)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return true, nil
}
