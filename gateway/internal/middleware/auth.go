package middleware

import (
	"net/http"
	"strings"

	"github.com/ManyLinesEditor/backend/gateway/internal/services"
	"github.com/gin-gonic/gin"
)

func Auth(tokens *services.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		userID, err := tokens.Verify(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Request.Header.Set("X-User-ID", userID.String())
		c.Next()
	}
}
