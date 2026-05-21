package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ManyLinesEditor/backend/storage/internal/models"
)

// FileRepo handles file metadata persistence in Postgres.
type FileRepo struct {
	db *pgxpool.Pool
}

// NewFileRepo creates a FileRepo with the provided connection pool.
func NewFileRepo(db *pgxpool.Pool) *FileRepo { return &FileRepo{db: db} }

// Create inserts a new file metadata row.
func (r *FileRepo) Create(ctx context.Context, f *models.File) error {
	const q = `
		INSERT INTO files (id, owner_id, name, content_type, size_bytes, bucket, object_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, q,
		f.ID, f.OwnerID, f.Name, f.ContentType, f.SizeBytes, f.Bucket, f.ObjectKey)
	return err
}

// FindByID returns file metadata by its UUID.
// Returns domain.ErrNotFound when no row is found.
func (r *FileRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.File, error) {
	const q = `
		SELECT id, owner_id, name, content_type, size_bytes, bucket, object_key, created_at, updated_at
		FROM files WHERE id = $1`

	var f models.File
	err := r.db.QueryRow(ctx, q, id).Scan(
		&f.ID, &f.OwnerID, &f.Name, &f.ContentType,
		&f.SizeBytes, &f.Bucket, &f.ObjectKey, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &f, nil
}
