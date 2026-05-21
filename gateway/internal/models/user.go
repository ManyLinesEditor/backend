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

// User represents an authenticated account.
// Login is a random UUID used as username.
type User struct {
	ID           uuid.UUID
	Login        uuid.UUID
	PasswordHash string
	CreatedAt    time.Time
}
