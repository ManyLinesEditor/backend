package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OwnerIDKey is the context key used to store the authenticated user ID.
const OwnerIDKey = "ownerID"

func GatewayAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		id, err := uuid.Parse(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad user id"})
			return
		}
		c.Set(OwnerIDKey, id)
		c.Next()
	}
}

// OwnerID extracts the authenticated user ID set by Auth.
// Panics if called outside a Auth-protected route.
func OwnerID(c *gin.Context) uuid.UUID {
	return c.MustGet(OwnerIDKey).(uuid.UUID)
}
