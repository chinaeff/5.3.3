package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware())

	router.GET("/test-secured-endpoint", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Authorized access"})
	})

	req := httptest.NewRequest("GET", "/test-secured-endpoint", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")

	req = httptest.NewRequest("GET", "/test-secured-endpoint", nil)
	req.Header.Set("Authorization", "valid_token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Authorized access")
}
