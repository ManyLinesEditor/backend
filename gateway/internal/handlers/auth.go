package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ManyLinesEditor/backend/gateway/internal/models"
	"github.com/ManyLinesEditor/backend/gateway/internal/services"
)

// AuthHandler provides HTTP handlers for auth endpoints.
type AuthHandler struct {
	auth *services.AuthService
}

// NewAuthHandler creates an AuthHandler with the provided AuthService.
func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register handles user registration.
//
// @Summary     User register (Login and Password is UUID)
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body services.Credentials true "User data"
// @Success     201 {object} models.TokenResponse
// @Failure     400 {object} models.ErrorResponse
// @Failure     409 {object} models.ErrorResponse "login already taken"
// @Router      /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var in services.Credentials
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.auth.Register(c.Request.Context(), in)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, models.ErrorResponse{Error: "login already taken"})
			return
		}
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.TokenResponse{Token: token})
}

// Login handles user authentication.
//
// @Summary     User login
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body services.Credentials true "User data"
// @Success     200 {object} models.TokenResponse
// @Failure     400 {object} models.ErrorResponse
// @Failure     401 {object} models.ErrorResponse
// @Router      /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var in services.Credentials
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.auth.Login(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{Token: token})
}
