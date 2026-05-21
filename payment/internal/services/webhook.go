package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ManyLinesEditor/backend/payment/internal/repositories"
	"github.com/ManyLinesEditor/backend/payment/internal/sse"
)

var ErrInvalidSignature = errors.New("invalid signature")

// Webhook body from acquiremock.
type AcquiremockPayload struct {
	PaymentID string  `json:"payment_id"`
	Reference string  `json:"reference"`
	Amount    float32 `json:"amount"`
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
}

type WebhookService struct {
	paymentRepo      repositories.PaymentRepo
	subscriptionRepo repositories.SubscriptionRepo
	broker           *sse.Broker
	webhookSecret    string
}

func NewWebhookService(
	paymentRepo repositories.PaymentRepo,
	subscriptionRepo repositories.SubscriptionRepo,
	broker *sse.Broker,
	webhookSecret string,
) *WebhookService {
	return &WebhookService{paymentRepo, subscriptionRepo, broker, webhookSecret}
}

// Handle processes the webhook: verifies the signature, updates the status,
// if the payment is successful, activates the subscription and pushes an SSE event.
func (s *WebhookService) Handle(ctx context.Context, body []byte, signature string) error {
	if !s.verifySignature(body, signature) {
		return ErrInvalidSignature
	}

	var payload AcquiremockPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("parse payload: %w", err)
	}

	payment, err := s.paymentRepo.UpdateStatus(ctx, payload.Reference, payload.Status)
	if err != nil {
		return nil
	}

	if payload.Status == "paid" {
		until := time.Now().AddDate(0, 1, 0)
		if err := s.subscriptionRepo.Upsert(ctx, payment.UserID, payment.FeatureID, until); err != nil {
			_ = err
		}
		s.publishFeatureAvailable(payment.UserID, payment.FeatureID, until)
	}

	return nil
}

func (s *WebhookService) publishFeatureAvailable(userID, featureID string, until time.Time) {
	event, _ := json.Marshal(map[string]string{
		"event":      "feature_available",
		"feature_id": featureID,
		"until":      until.Format(time.RFC3339),
	})
	s.broker.Publish(userID, string(event))
}

func (s *WebhookService) verifySignature(body []byte, signature string) bool {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()

	var m map[string]json.RawMessage
	if err := dec.Decode(&m); err != nil {
		return false
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteString(", ")
		}
		keyJSON, _ := json.Marshal(k)
		b.Write(keyJSON)
		b.WriteString(": ")
		b.Write(m[k])
	}
	b.WriteByte('}')

	canonical := b.String()

	mac := hmac.New(sha256.New, []byte(s.webhookSecret))
	mac.Write([]byte(canonical))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}
