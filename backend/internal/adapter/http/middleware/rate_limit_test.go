package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestUserRateLimiterLimitsAndResetsPerWindow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Date(2026, 9, 27, 12, 0, 0, 0, time.UTC)
	limiter := NewUserRateLimiter(2, time.Minute)
	limiter.now = func() time.Time { return now }

	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set(authenticatedUserIDKey, "user-1")
	})
	router.POST("/ai", limiter.Handle(), func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })

	for requestNumber := 1; requestNumber <= 2; requestNumber++ {
		responseRecorder := httptest.NewRecorder()
		router.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodPost, "/ai", nil))
		if responseRecorder.Code != http.StatusNoContent {
			t.Fatalf("request %d status = %d, want 204", requestNumber, responseRecorder.Code)
		}
	}

	limitedResponse := httptest.NewRecorder()
	router.ServeHTTP(limitedResponse, httptest.NewRequest(http.MethodPost, "/ai", nil))
	if limitedResponse.Code != http.StatusTooManyRequests || limitedResponse.Header().Get("Retry-After") == "" {
		t.Fatalf("limited response status = %d, retry-after = %q", limitedResponse.Code, limitedResponse.Header().Get("Retry-After"))
	}

	now = now.Add(time.Minute)
	resetResponse := httptest.NewRecorder()
	router.ServeHTTP(resetResponse, httptest.NewRequest(http.MethodPost, "/ai", nil))
	if resetResponse.Code != http.StatusNoContent {
		t.Fatalf("request after reset status = %d, want 204", resetResponse.Code)
	}
}
