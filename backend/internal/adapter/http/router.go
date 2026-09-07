package http

import (
	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/handler"
	"job-copilot-backend/internal/adapter/http/middleware"
	analysisapp "job-copilot-backend/internal/application/analysis"
	"job-copilot-backend/internal/port"
)

type Dependencies struct {
	AuthProvider     port.AuthProvider
	AnalyzeJDService *analysisapp.AnalyzeJDService
}

func NewRouter(frontendOrigin string, dependencies Dependencies) *gin.Engine {
	engine := gin.Default()

	engine.Use(middleware.CORS(frontendOrigin))

	engine.GET("/health", handler.Health)

	apiV1 := engine.Group("/api/v1")
	apiV1.GET("/health", handler.Health)

	aiRoutes := apiV1.Group("/ai")
	aiRoutes.Use(middleware.Authenticate(dependencies.AuthProvider))
	aiRoutes.POST("/analyze-jd", handler.NewAnalyzeJDHandler(dependencies.AnalyzeJDService).Handle)

	engine.NoRoute(handler.NotFound)

	return engine
}
