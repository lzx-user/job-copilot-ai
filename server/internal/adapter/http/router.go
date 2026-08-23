package http

import (
	"github.com/gin-gonic/gin"

	"job-copilot-server/internal/adapter/http/handler"
	"job-copilot-server/internal/adapter/http/middleware"
)

func NewRouter(frontendOrigin string) *gin.Engine {
	engine := gin.Default()

	engine.Use(middleware.CORS(frontendOrigin))

	engine.GET("/health", handler.Health)
	engine.NoRoute(handler.NotFound)

	return engine
}
