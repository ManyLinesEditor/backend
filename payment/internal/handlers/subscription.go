package handlers

import (
	"net/http"

	"github.com/ManyLinesEditor/backend/payment/internal/middleware"
	"github.com/ManyLinesEditor/backend/payment/internal/models"
	"github.com/ManyLinesEditor/backend/payment/internal/services"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subscriptionSvc *services.SubscriptionService
}

func NewSubscriptionHandler(svc *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc}
}

// ListActive returns all active subscriptions for the authenticated user.
//
// @Summary     List active subscriptions
// @Description Returns all active subscriptions for the authenticated user
// @Tags        subscriptions
// @Produce     json
// @Success     200 {array}  models.Subscription
// @Failure     500 {object} models.ErrorResponse
// @Security    BearerAuth
// @Router      /subscriptions [get]
func (h *SubscriptionHandler) ListActive(c *gin.Context) {
	userID := c.GetString(middleware.OwnerIDKey)

	subs, err := h.subscriptionSvc.ListActive(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal error"})
		return
	}
	if subs == nil {
		subs = make([]*models.Subscription, 0)
	}

	c.JSON(http.StatusOK, subs)
}
