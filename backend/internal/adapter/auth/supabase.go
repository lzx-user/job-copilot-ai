package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"job-copilot-backend/internal/port"
)

var (
	ErrInvalidAuthConfig = errors.New("invalid auth config")
)

// SupabaseAuthAdapter 使用 Supabase Auth 实现核心层要求的认证能力。
type SupabaseAuthAdapter struct {
	baseURL    string
	anonKey    string
	httpClient *http.Client
}

// 编译器会在方法签名不匹配时直接报错，明确保证该 Adapter 实现了 AuthProvider。
var _ port.AuthProvider = (*SupabaseAuthAdapter)(nil)

func NewSupabaseAuthAdapter(baseURL string, anonKey string, httpClient *http.Client) (*SupabaseAuthAdapter, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	anonKey = strings.TrimSpace(anonKey)
	if baseURL == "" || anonKey == "" {
		return nil, ErrInvalidAuthConfig
	}

	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &SupabaseAuthAdapter{
		baseURL:    baseURL,
		anonKey:    anonKey,
		httpClient: httpClient,
	}, nil
}

func (adapter *SupabaseAuthAdapter) VerifyAccessToken(
	ctx context.Context,
	accessToken string,
) (port.AuthenticatedUser, error) {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return port.AuthenticatedUser{}, port.ErrUnauthenticated
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		adapter.baseURL+"/auth/v1/user",
		nil,
	)
	if err != nil {
		return port.AuthenticatedUser{}, fmt.Errorf("create supabase auth request: %w", err)
	}

	request.Header.Set("apikey", adapter.anonKey)
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := adapter.httpClient.Do(request)
	if err != nil {
		return port.AuthenticatedUser{}, fmt.Errorf("verify access token: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return port.AuthenticatedUser{}, port.ErrUnauthenticated
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return port.AuthenticatedUser{}, fmt.Errorf("supabase auth returned status %d", response.StatusCode)
	}

	var authUser struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&authUser); err != nil {
		return port.AuthenticatedUser{}, fmt.Errorf("decode supabase auth response: %w", err)
	}
	if strings.TrimSpace(authUser.ID) == "" {
		return port.AuthenticatedUser{}, port.ErrUnauthenticated
	}

	return port.AuthenticatedUser{UserID: authUser.ID}, nil
}
