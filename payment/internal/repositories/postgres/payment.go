package postgres

import (
	"context"
	"fmt"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepo struct {
	pool *pgxpool.Pool
}

func NewPaymentRepo(pool *pgxpool.Pool) *PaymentRepo {
	return &PaymentRepo{pool: pool}
}

func (r *PaymentRepo) Create(ctx context.Context, userID, featureID string, amount int, reference string) (*models.Payment, error) {
	const q = `
		INSERT INTO payments (user_id, feature_id, amount, reference)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, feature_id, status, amount, created_at, updated_at`

	p := &models.Payment{}
	err := r.pool.QueryRow(ctx, q, userID, featureID, amount, reference).
		Scan(&p.ID, &p.UserID, &p.FeatureID, &p.Status, &p.Amount, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("payment create: %w", err)
	}
	return p, nil
}

func (r *PaymentRepo) SetCheckoutURL(ctx context.Context, id, url string) error {
	const q = `UPDATE payments SET checkout_url = $1, updated_at = now() WHERE id = $2`
	_, err := r.pool.Exec(ctx, q, url, id)
	return err
}

// UpdateStatus updates the status and returns the payment.
func (r *PaymentRepo) UpdateStatus(ctx context.Context, id, status string) (*models.Payment, error) {
	const q = `
		UPDATE payments SET status = $1, updated_at = now()
		WHERE id = $2
		RETURNING id, user_id, feature_id, status, amount, created_at, updated_at`

	p := &models.Payment{}
	err := r.pool.QueryRow(ctx, q, status, id).
		Scan(&p.ID, &p.UserID, &p.FeatureID, &p.Status, &p.Amount, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("payment update status: %w", err)
	}
	return p, nil
}

func (r *PaymentRepo) ListByUser(ctx context.Context, userID string) ([]*models.Payment, error) {
	const q = `
		SELECT id, user_id, feature_id, status, amount, COALESCE(checkout_url,''), created_at, updated_at
		FROM payments WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Payment
	for rows.Next() {
		p := &models.Payment{}
		if err := rows.Scan(&p.ID, &p.UserID, &p.FeatureID, &p.Status,
			&p.Amount, &p.CheckoutURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}
