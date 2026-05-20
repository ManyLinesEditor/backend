package models

import "github.com/google/uuid"

type ErrorResponse struct {
	Error string `json:"error"`
}

// UploadResult is returned after a successful upload.
type UploadResult struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
