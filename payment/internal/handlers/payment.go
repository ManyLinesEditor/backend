package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/ManyLinesEditor/backend/payment/internal/middleware"
	"github.com/ManyLinesEditor/backend/payment/internal/models"
	"github.com/ManyLinesEditor/backend/payment/internal/services"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentSvc *services.PaymentService
}

func NewPaymentHandler(paymentSvc *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentSvc}
}

// POST /payments
type createPaymentRequest struct {
	FeatureID string `json:"feature_id" binding:"required"`
	Amount    int    `json:"amount"     binding:"required,min=1"`
}

type createPaymentResponse struct {
	PaymentID   string `json:"payment_id"`
	CheckoutURL string `json:"checkout_url"`
}

// Create initiates a new payment.
//
// @Summary     Create payment
// @Description Initiates a new payment for the specified feature
// @Tags        payments
// @Accept      json
// @Produce     json
// @Param       request body createPaymentRequest true "Payment data"
// @Success     201 {object} createPaymentResponse
// @Failure     400 {object} models.ErrorResponse
// @Failure     502 {object} models.ErrorResponse "payment gateway unavailable"
// @Failure     500 {object} models.ErrorResponse
// @Security    BearerAuth
// @Router      /payments [post]
func (h *PaymentHandler) Create(c *gin.Context) {
	userID := c.GetString(middleware.OwnerIDKey)

	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.paymentSvc.Create(c.Request.Context(), userID, req.FeatureID, req.Amount)
	if err != nil {
		if errors.Is(err, services.ErrAcquiremockUnavailable) {
			c.JSON(http.StatusBadGateway, models.ErrorResponse{Error: "payment gateway unavailable"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal error"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, createPaymentResponse{
		PaymentID:   result.Payment.ID,
		CheckoutURL: result.CheckoutURL,
	})
}

// List returns all payments for the authenticated user.
//
// @Summary     List payments
// @Description Returns all payments for the authenticated user
// @Tags        payments
// @Produce     json
// @Success     200 {array} models.Payment
// @Failure     500 {object} models.ErrorResponse
// @Security    BearerAuth
// @Router      /payments [get]
func (h *PaymentHandler) List(c *gin.Context) {
	userID := c.GetString(middleware.OwnerIDKey)

	payments, err := h.paymentSvc.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal error"})
		return
	}

	c.JSON(http.StatusOK, payments)
}
