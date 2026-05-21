package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type WebhookResponse struct {
	Status string `json:"status"`
}
