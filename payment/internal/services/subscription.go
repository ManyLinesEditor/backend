package services

import (
	"context"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
	"github.com/ManyLinesEditor/backend/payment/internal/repositories"
)

type SubscriptionService struct {
	subscriptionRepo repositories.SubscriptionRepo
}

func NewSubscriptionService(repo repositories.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo}
}

func (s *SubscriptionService) ListActive(ctx context.Context, userID string) ([]*models.Subscription, error) {
	return s.subscriptionRepo.ListActive(ctx, userID)
}

func (s *SubscriptionService) HasActive(ctx context.Context, userID, featureID string) (bool, error) {
	return s.subscriptionRepo.HasActive(ctx, userID, featureID)
}
