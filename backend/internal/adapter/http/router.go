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
	AnalysisHistoryService      *analysisapp.HistoryService
	StartInterviewService       *interviewapp.StartInterviewService
	ListInterviewOptionsService *interviewapp.ListInterviewOptionsService
	InterviewTurnService        *interviewapp.TurnService
	InterviewSessionService     *interviewapp.SessionService
	InterviewReportService      *interviewapp.ReportService
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
	analysisHistoryHandler := handler.NewJDAnalysisHistoryHandler(dependencies.AnalysisHistoryService)
	aiRoutes.GET("/jd-analyses", analysisHistoryHandler.List)
	aiRoutes.GET("/jd-analyses/:id", analysisHistoryHandler.Get)
	aiRoutes.GET("/interview/options", handler.NewListInterviewOptionsHandler(dependencies.ListInterviewOptionsService).Handle)
	aiRoutes.POST("/interview/start", handler.NewStartInterviewHandler(dependencies.StartInterviewService).Handle)
	aiRoutes.POST("/interview/turn", handler.NewInterviewTurnHandler(dependencies.InterviewTurnService).Handle)
	aiRoutes.POST("/interview/report", handler.NewInterviewReportHandler(dependencies.InterviewReportService).Handle)
	aiRoutes.GET("/interview/sessions/:id", handler.NewInterviewSessionHandler(dependencies.InterviewSessionService).Get)

	engine.NoRoute(handler.NotFound)

	return engine
}
