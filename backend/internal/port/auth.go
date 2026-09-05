package port

import "context"

// AuthenticatedUser 是访问令牌验证成功后得到的可信用户身份。
type AuthenticatedUser struct {
	UserID string
}

// AuthProvider 定义核心业务需要的最小认证能力。
type AuthProvider interface {
	VerifyAccessToken(ctx context.Context, accessToken string) (AuthenticatedUser, error)
}
