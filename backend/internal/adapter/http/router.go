package http

import (
	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/handler"
	"job-copilot-backend/internal/adapter/http/middleware"
	analysisapp "job-copilot-backend/internal/application/analysis"
	interviewapp "job-copilot-backend/internal/application/interview"
	"job-copilot-backend/internal/port"
)

type Dependencies struct {
	AuthProvider                port.AuthProvider
	AnalyzeJDService            *analysisapp.AnalyzeJDService
	StartInterviewService       *interviewapp.StartInterviewService
	ListInterviewOptionsService *interviewapp.ListInterviewOptionsService
	InterviewTurnService        *interviewapp.TurnService
	InterviewSessionService     *interviewapp.SessionService
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
	aiRoutes.GET("/interview/options", handler.NewListInterviewOptionsHandler(dependencies.ListInterviewOptionsService).Handle)
	aiRoutes.POST("/interview/start", handler.NewStartInterviewHandler(dependencies.StartInterviewService).Handle)
	aiRoutes.POST("/interview/turn", handler.NewInterviewTurnHandler(dependencies.InterviewTurnService).Handle)
	aiRoutes.GET("/interview/sessions/:id", handler.NewInterviewSessionHandler(dependencies.InterviewSessionService).Get)

	engine.NoRoute(handler.NotFound)

	return engine
}
