package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
	"github.com/ManyLinesEditor/backend/payment/internal/services"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	webhookSvc *services.WebhookService
}

func NewWebhookHandler(webhookSvc *services.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookSvc}
}

// Handle processes a webhook from AcquireMock payment gateway.
//
// @Summary     AcquireMock webhook
// @Description Receives payment confirmation events from AcquireMock.
// @Description Validates HMAC signature from X-Signature header.
// @Tags        webhooks
// @Accept      application/octet-stream
// @Produce     json
// @Param       X-Signature header string true "HMAC signature"
// @Success     200 {object} models.WebhookResponse
// @Failure     400 {object} models.ErrorResponse
// @Failure     403 {object} models.ErrorResponse "invalid signature"
// @Router      /webhooks/acquiremock [post]
func (h *WebhookHandler) Handle(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "can't read body"})
		return
	}

	if err := h.webhookSvc.Handle(c.Request.Context(), body, c.GetHeader("X-Signature")); err != nil {
		if errors.Is(err, services.ErrInvalidSignature) {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "invalid signature"})
			return
		}
		log.Println(err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.WebhookResponse{Status: "ok"})
}
