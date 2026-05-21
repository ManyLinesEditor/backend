package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ManyLinesEditor/backend/storage/internal/models"
)

// FileRepository defines file metadata persistence behavior.
type FileRepo interface {
	Create(ctx context.Context, f *models.File) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.File, error)
}
