package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/port"
	"job-copilot-backend/pkg/response"
)

const authenticatedUserIDKey = "authenticatedUserID"

func Authenticate(authProvider port.AuthProvider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if authProvider == nil {
			response.Error(ctx, http.StatusServiceUnavailable, "AUTH_NOT_CONFIGURED", "认证服务未配置")
			ctx.Abort()
			return
		}

		authorization := strings.TrimSpace(ctx.GetHeader("Authorization"))
		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
			ctx.Abort()
			return
		}

		user, err := authProvider.VerifyAccessToken(ctx.Request.Context(), parts[1])
		if err != nil {
			if errors.Is(err, port.ErrUnauthenticated) {
				response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "登录状态已失效")
			} else {
				response.Error(ctx, http.StatusBadGateway, "AUTH_UPSTREAM_ERROR", "认证服务暂时不可用")
			}
			ctx.Abort()
			return
		}

		ctx.Set(authenticatedUserIDKey, user.UserID)
		ctx.Next()
	}
}

func AuthenticatedUserID(ctx *gin.Context) (string, bool) {
	userID, exists := ctx.Get(authenticatedUserIDKey)
	if !exists {
		return "", false
	}
	value, ok := userID.(string)
	return value, ok && strings.TrimSpace(value) != ""
}
