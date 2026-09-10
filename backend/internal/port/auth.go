package port

import (
	"context"
	"errors"
)

var ErrUnauthenticated = errors.New("unauthenticated")

type accessTokenContextKey struct{}

// AuthenticatedUser 是访问令牌验证成功后得到的可信用户身份。
type AuthenticatedUser struct {
	UserID string
}

// AuthProvider 定义核心业务需要的最小认证能力。
type AuthProvider interface {
	VerifyAccessToken(ctx context.Context, accessToken string) (AuthenticatedUser, error)
}

// WithAuthenticatedAccessToken 只应在令牌校验成功后调用，供需要遵守 RLS 的外部 Adapter 使用。
func WithAuthenticatedAccessToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, accessTokenContextKey{}, accessToken)
}

func AuthenticatedAccessToken(ctx context.Context) (string, bool) {
	accessToken, ok := ctx.Value(accessTokenContextKey{}).(string)
	return accessToken, ok && accessToken != ""
}
