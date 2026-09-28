package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/pkg/response"
)

type rateLimitEntry struct {
	windowStarted time.Time
	count         int
}

// UserRateLimiter 对已认证用户的高成本 AI 请求做单实例固定窗口限流。
type UserRateLimiter struct {
	mu      sync.Mutex
	entries map[string]rateLimitEntry
	max     int
	window  time.Duration
	now     func() time.Time
}

func NewUserRateLimiter(max int, window time.Duration) *UserRateLimiter {
	if max <= 0 {
		max = 10
	}
	if window <= 0 {
		window = time.Minute
	}
	return &UserRateLimiter{
		entries: make(map[string]rateLimitEntry),
		max:     max, window: window, now: time.Now,
	}
}

func (limiter *UserRateLimiter) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := AuthenticatedUserID(ctx)
		if !ok {
			response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
			ctx.Abort()
			return
		}

		now := limiter.now()
		limiter.mu.Lock()
		entry := limiter.entries[userID]
		if entry.windowStarted.IsZero() || now.Sub(entry.windowStarted) >= limiter.window {
			entry = rateLimitEntry{windowStarted: now}
		}
		if entry.count >= limiter.max {
			retryAfter := int((limiter.window - now.Sub(entry.windowStarted)).Seconds()) + 1
			limiter.mu.Unlock()
			ctx.Header("Retry-After", strconv.Itoa(retryAfter))
			response.Error(ctx, http.StatusTooManyRequests, "RATE_LIMITED", "AI 请求过于频繁，请稍后再试")
			ctx.Abort()
			return
		}
		entry.count++
		limiter.entries[userID] = entry
		if len(limiter.entries) > 10000 {
			for key, value := range limiter.entries {
				if now.Sub(value.windowStarted) >= limiter.window {
					delete(limiter.entries, key)
				}
			}
		}
		limiter.mu.Unlock()
		ctx.Next()
	}
}
