package handlers

import (
	"io"
	"net/http"

	"github.com/ManyLinesEditor/backend/payment/internal/middleware"
	"github.com/ManyLinesEditor/backend/payment/internal/sse"
	"github.com/gin-gonic/gin"
)

type SSEHandler struct {
	broker *sse.Broker
}

func NewSSEHandler(broker *sse.Broker) *SSEHandler {
	return &SSEHandler{broker}
}

// Stream opens a Server-Sent Events connection for the authenticated user.
// The connection stays open until the client disconnects.
// When a payment is confirmed, a "feature" event is sent with a JSON payload.
//
// @Summary     Subscribe to feature events
// @Description Opens a persistent SSE connection. Sends a "ping" event on connect,
// @Description then "feature" events when payments are confirmed.
// @Tags        sse
// @Produce     text/event-stream
// @Success     200 {string} string "event: feature\ndata: {payload}"
// @Failure     401 {object} models.ErrorResponse
// @Security    BearerAuth
// @Router      /sse/features [get]
func (h *SSEHandler) Stream(c *gin.Context) {
	userID := c.GetString(middleware.OwnerIDKey)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ch := h.broker.Subscribe(userID)
	defer h.broker.Unsubscribe(userID, ch)

	c.SSEvent("ping", "connected")
	c.Writer.Flush()

	c.Stream(func(w io.Writer) bool {
		select {
		case payload, ok := <-ch:
			if !ok {
				return false
			}
			c.SSEvent("feature", payload)
			return true

		case <-c.Request.Context().Done():
			return false
		}
	})

	c.Status(http.StatusOK)
}
