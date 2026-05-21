package repositories

import (
	"context"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
)

type PaymentRepo interface {
	Create(ctx context.Context, userID, featureID string, amount int, reference string) (*models.Payment, error)
	SetCheckoutURL(ctx context.Context, id, url string) error
	UpdateStatus(ctx context.Context, id, status string) (*models.Payment, error)
	ListByUser(ctx context.Context, userID string) ([]*models.Payment, error)
}
