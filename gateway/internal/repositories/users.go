package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/ManyLinesEditor/backend/gateway/internal/models"
)

// UserRepository defines user persistence behavior.
type UserRepo interface {
	Create(ctx context.Context, login, passwordHash string) (*models.User, error)
	FindByLogin(ctx context.Context, login uuid.UUID) (*models.User, error)
}
