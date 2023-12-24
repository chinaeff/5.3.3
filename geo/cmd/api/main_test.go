package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkMoveCourierHandler(b *testing.B) {
	router := gin.Default()

	router.POST("/move-courier", func(c *gin.Context) {
		c.Request, _ = http.NewRequest("POST", "/move-courier", nil)

		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = http.NoBody

		router.ServeHTTP(c.Writer, c.Request)
	})

	req, _ := http.NewRequest("POST", "/move-courier", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Создаем HTTP запрос
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}
