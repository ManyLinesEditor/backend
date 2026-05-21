package services

import (
	"context"
	"errors"

	"github.com/ManyLinesEditor/backend/gateway/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrAlreadyExists = errors.New("user already exists")
var ErrBadCredentials = errors.New("invalid credentials")

type AuthService struct {
	users  repositories.UserRepo
	tokens *TokenService
}

func NewAuthService(users repositories.UserRepo, tokens *TokenService) *AuthService {
	return &AuthService{users, tokens}
}

type Credentials struct {
	Login    string `json:"login"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthService) Register(ctx context.Context, c Credentials) (string, error) {
	loginUUID, err := uuid.Parse(c.Login)
	if err != nil {
		return "", errors.New("login must be a valid UUID")
	}
	if _, err = uuid.Parse(c.Password); err != nil {
		return "", errors.New("password must be a valid UUID")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user, err := s.users.Create(ctx, loginUUID.String(), string(hash))
	if err != nil {
		return "", ErrAlreadyExists
	}

	return s.tokens.Issue(user.ID)
}

func (s *AuthService) Login(ctx context.Context, c Credentials) (string, error) {
	loginUUID, err := uuid.Parse(c.Login)
	if err != nil {
		return "", ErrBadCredentials
	}

	user, err := s.users.FindByLogin(ctx, loginUUID)
	if err != nil {
		return "", ErrBadCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(c.Password)); err != nil {
		return "", ErrBadCredentials
	}

	return s.tokens.Issue(user.ID)
}
