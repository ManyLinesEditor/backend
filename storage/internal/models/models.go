// Package domain contains core entity types and shared errors.
package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Sentinel errors returned by repositories.
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

// File holds object metadata; the binary is stored in MinIO.
type File struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Name        string
	ContentType string
	SizeBytes   int64
	Bucket      string
	ObjectKey   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
