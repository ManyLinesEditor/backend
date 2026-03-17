package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ManyLines/internal/handlers"

	"github.com/gin-gonic/gin"
)

func TestDefaultSmoke(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", handlers.Default)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
