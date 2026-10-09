package http

import (
	"time"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/handler"
	"job-copilot-backend/internal/adapter/http/middleware"
	analysisapp "job-copilot-backend/internal/application/analysis"
	careerapp "job-copilot-backend/internal/application/career"
	dashboardapp "job-copilot-backend/internal/application/dashboard"
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
	InterviewHistoryService     *interviewapp.HistoryService
	AIRateLimit                 int
	DashboardService            *dashboardapp.Service
	CareerService               *careerapp.Service
}

func NewRouter(frontendOrigin string, dependencies Dependencies) *gin.Engine {
	engine := gin.Default()

	engine.Use(middleware.CORS(frontendOrigin))

	engine.GET("/health", handler.Health)

	apiV1 := engine.Group("/api/v1")
	apiV1.GET("/health", handler.Health)

	aiRoutes := apiV1.Group("/ai")
	aiRoutes.Use(middleware.Authenticate(dependencies.AuthProvider))
	aiRoutes.GET("/dashboard", handler.NewDashboardHandler(dependencies.DashboardService).Get)
	careerHandler := handler.NewCareerHandler(dependencies.CareerService)
	careerRoutes := apiV1.Group("/career")
	careerRoutes.Use(middleware.Authenticate(dependencies.AuthProvider))
	careerRoutes.GET("/workspace", careerHandler.GetWorkspace)
	careerRoutes.POST("/resumes", careerHandler.CreateResume)
	careerRoutes.POST("/resumes/:id/versions", careerHandler.CreateResumeVersion)
	careerRoutes.POST("/jobs", careerHandler.CreateJob)
	careerRoutes.POST("/jobs/:id/jd-versions", careerHandler.CreateJDVersion)
	careerRoutes.POST("/applications", careerHandler.CreateApplication)
	careerRoutes.POST("/events", careerHandler.AddEvent)
	careerRoutes.POST("/interviews", careerHandler.CreateInterview)
	careerRoutes.PATCH("/interviews/:id", careerHandler.UpdateInterview)
	careerRoutes.PUT("/retrospectives", careerHandler.SaveRetrospective)
	careerRoutes.PUT("/offers", careerHandler.SaveOffer)
	aiRateLimit := middleware.NewUserRateLimiter(dependencies.AIRateLimit, time.Minute).Handle()
	aiRoutes.POST("/analyze-jd", aiRateLimit, handler.NewAnalyzeJDHandler(dependencies.AnalyzeJDService).Handle)
	analysisHistoryHandler := handler.NewJDAnalysisHistoryHandler(dependencies.AnalysisHistoryService)
	aiRoutes.GET("/jd-analyses", analysisHistoryHandler.List)
	aiRoutes.GET("/jd-analyses/:id", analysisHistoryHandler.Get)
	aiRoutes.DELETE("/jd-analyses/:id", analysisHistoryHandler.Delete)
	aiRoutes.GET("/interview/options", handler.NewListInterviewOptionsHandler(dependencies.ListInterviewOptionsService).Handle)
	aiRoutes.POST("/interview/start", aiRateLimit, handler.NewStartInterviewHandler(dependencies.StartInterviewService).Handle)
	aiRoutes.POST("/interview/turn", aiRateLimit, handler.NewInterviewTurnHandler(dependencies.InterviewTurnService).Handle)
	aiRoutes.POST("/interview/report", aiRateLimit, handler.NewInterviewReportHandler(dependencies.InterviewReportService).Handle)
	interviewHistoryHandler := handler.NewInterviewHistoryHandler(dependencies.InterviewHistoryService)
	aiRoutes.GET("/interview/sessions", interviewHistoryHandler.List)
	aiRoutes.GET("/interview/sessions/:id", handler.NewInterviewSessionHandler(dependencies.InterviewSessionService).Get)
	aiRoutes.DELETE("/interview/sessions/:id", interviewHistoryHandler.Delete)

	engine.NoRoute(handler.NotFound)

	return engine
}
