package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/ManyLinesEditor/backend/gateway/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) FindByLogin(ctx context.Context, login uuid.UUID) (*models.User, error) {
	const q = `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	row := r.pool.QueryRow(ctx, q, login)

	var u models.User
	if err := row.Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// Register creates a new user (for testing/registration).
// The password_hash is passed already hashed from the outside.
func (r *UserRepo) Create(ctx context.Context, login, passwordHash string) (*models.User, error) {
	const q = `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
		RETURNING id, login, password_hash, created_at`
	u := &models.User{}
	err := r.pool.QueryRow(ctx, q, login, passwordHash).
		Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user register: %w", err)
	}
	return u, nil
}
