package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
	"github.com/ManyLinesEditor/backend/payment/internal/repositories"
	"github.com/google/uuid"
)

var ErrAcquiremockUnavailable = errors.New("acquiremock unavailable")

type PaymentService struct {
	paymentRepo repositories.PaymentRepo
	acquireURL  string
	baseURL     string
}

func NewPaymentService(
	paymentRepo repositories.PaymentRepo,
	acquireURL, baseURL string,
) *PaymentService {
	return &PaymentService{paymentRepo, acquireURL, baseURL}
}

type CreatePaymentResult struct {
	Payment     *models.Payment
	CheckoutURL string
}

func (s *PaymentService) Create(ctx context.Context, userID, featureID string, amount int) (*CreatePaymentResult, error) {
	reference := uuid.New().String()
	payment, err := s.paymentRepo.Create(ctx, userID, featureID, amount, reference)
	if err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	checkoutURL, err := s.createInvoice(payment)
	if err != nil {
		return nil, err
	}

	if err := s.paymentRepo.SetCheckoutURL(ctx, payment.ID, checkoutURL); err != nil {
		_ = err
	}
	payment.CheckoutURL = checkoutURL

	return &CreatePaymentResult{Payment: payment, CheckoutURL: checkoutURL}, nil
}

func (s *PaymentService) ListByUser(ctx context.Context, userID string) ([]*models.Payment, error) {
	return s.paymentRepo.ListByUser(ctx, userID)
}

func (s *PaymentService) createInvoice(p *models.Payment) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"amount":      p.Amount,
		"reference":   p.ID,
		"webhookUrl":  fmt.Sprintf("%s/webhooks/acquiremock", s.baseURL),
		"redirectUrl": fmt.Sprintf("%s/payments/%s/result", s.baseURL, p.ID),
	})

	resp, err := http.Post(
		fmt.Sprintf("%s/api/create-invoice", s.acquireURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil || resp.StatusCode >= 400 {
		return "", ErrAcquiremockUnavailable
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()
	var result struct {
		PageURL string `json:"pageUrl"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil || result.PageURL == "" {
		return "", fmt.Errorf("bad acquiremock response: %w", err)
	}
	return result.PageURL, err
}
