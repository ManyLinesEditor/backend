// Package handlers includes HTTP-handlers
package handlers

import (
	"github.com/gin-gonic/gin"
)

// Default handles requests to the root route and returns a welcome message in JSON format.
func Default(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World",
	})
}
