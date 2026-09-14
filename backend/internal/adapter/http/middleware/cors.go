package middleware

import (
	"net/url"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigin string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: localOriginAliases(allowedOrigin),
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
	})
}

func localOriginAliases(configuredOrigin string) []string {
	origins := make([]string, 0, 3)
	seen := make(map[string]struct{})
	add := func(origin string) {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			return
		}
		if _, exists := seen[origin]; exists {
			return
		}
		seen[origin] = struct{}{}
		origins = append(origins, origin)
	}

	for _, configured := range strings.Split(configuredOrigin, ",") {
		configured = strings.TrimSpace(configured)
		add(configured)

		parsed, err := url.Parse(configured)
		if err != nil || !isLoopbackHost(parsed.Hostname()) {
			continue
		}

		port := parsed.Port()
		if port != "" {
			port = ":" + port
		}
		// 浏览器可能把 localhost 解析为 IPv4 或 IPv6，这三种写法在本地开发中等价。
		add(parsed.Scheme + "://localhost" + port)
		add(parsed.Scheme + "://127.0.0.1" + port)
		add(parsed.Scheme + "://[::1]" + port)
	}

	return origins
}

func isLoopbackHost(host string) bool {
	switch strings.ToLower(host) {
	case "localhost", "127.0.0.1", "::1":
		return true
	default:
		return false
	}
}
