package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Device struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Fingerprint string    `json:"fingerprint"`
	Revoked     bool      `json:"revoked"`
	CreatedAt   time.Time `json:"created_at"`
}

type Feature struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Subscription struct {
	UserID    string    `json:"user_id"`
	FeatureID string    `json:"feature_id"`
	Until     time.Time `json:"until"`
}

type Payment struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	FeatureID   string    `json:"feature_id"`
	Status      string    `json:"status"`
	Amount      int       `json:"amount"`
	CheckoutURL string    `json:"checkout_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
