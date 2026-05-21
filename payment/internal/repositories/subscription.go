package repositories

import (
	"context"
	"time"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
)

type SubscriptionRepo interface {
	Upsert(ctx context.Context, userID, featureID string, until time.Time) error
	ListActive(ctx context.Context, userID string) ([]*models.Subscription, error)
	HasActive(ctx context.Context, userID, featureID string) (bool, error)
}
