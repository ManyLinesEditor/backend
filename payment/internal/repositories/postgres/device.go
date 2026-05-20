package postgres

import (
	"context"
	"fmt"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepo struct {
	pool *pgxpool.Pool
}

func NewDeviceRepo(pool *pgxpool.Pool) *DeviceRepo {
	return &DeviceRepo{pool: pool}
}

// FindOrCreate returns the device or creates a new one.
// If the device is revoked, it returns it "as is." The handler will check the flag.
func (r *DeviceRepo) FindOrCreate(ctx context.Context, userID, fingerprint string) (*models.Device, error) {
	const q = `
		INSERT INTO devices (user_id, fingerprint)
		VALUES ($1, $2)
		ON CONFLICT (user_id, fingerprint) DO UPDATE SET user_id = EXCLUDED.user_id
		RETURNING id, user_id, fingerprint, revoked, created_at`

	d := &models.Device{}
	err := r.pool.QueryRow(ctx, q, userID, fingerprint).
		Scan(&d.ID, &d.UserID, &d.Fingerprint, &d.Revoked, &d.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("device upsert: %w", err)
	}
	return d, nil
}

func (r *DeviceRepo) GetByID(ctx context.Context, id string) (*models.Device, error) {
	const q = `SELECT id, user_id, fingerprint, revoked, created_at FROM devices WHERE id = $1`
	d := &models.Device{}
	err := r.pool.QueryRow(ctx, q, id).
		Scan(&d.ID, &d.UserID, &d.Fingerprint, &d.Revoked, &d.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("device get: %w", err)
	}
	return d, nil
}
