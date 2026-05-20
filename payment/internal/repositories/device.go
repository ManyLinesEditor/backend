package repositories

import (
	"context"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
)

type DeviceRepo interface {
	FindOrCreate(ctx context.Context, userID, fingerprint string) (*models.Device, error)
	GetByID(ctx context.Context, id string) (*models.Device, error)
}
