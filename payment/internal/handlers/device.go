package handlers

import (
	"net/http"

	"github.com/ManyLinesEditor/backend/payment/internal/models"
	"github.com/ManyLinesEditor/backend/payment/internal/repositories"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceRepo repositories.DeviceRepo
}

func NewDeviceHandler(repo repositories.DeviceRepo) *DeviceHandler {
	return &DeviceHandler{repo}
}

type registerDeviceRequest struct {
	Fingerprint string `json:"fingerprint" binding:"required"`
}

type registerDeviceResponse struct {
	DeviceSyncID string `json:"device_sync_id"`
}

// Register registers a device for the authenticated user.
//
// @Summary     Register device
// @Description Called once after login. Returns device_sync_id used as device_id in the deltas table.
// @Tags        devices
// @Accept      json
// @Produce     json
// @Param       request body registerDeviceRequest true "Device data"
// @Success     200 {object} registerDeviceResponse
// @Failure     400 {object} models.ErrorResponse
// @Failure     403 {object} models.ErrorResponse "device revoked"
// @Failure     500 {object} models.ErrorResponse
// @Router      /devices [post]
func (h *DeviceHandler) Register(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")

	var req registerDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	device, err := h.deviceRepo.FindOrCreate(c.Request.Context(), userID, req.Fingerprint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal error"})
		return
	}
	if device.Revoked {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "device revoked"})
		return
	}

	c.JSON(http.StatusOK, registerDeviceResponse{DeviceSyncID: device.ID})
}
