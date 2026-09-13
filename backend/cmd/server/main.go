package main

import (
	"log"
	"net/http"

	"job-copilot-backend/internal/adapter/ai"
	"job-copilot-backend/internal/adapter/auth"
	httpadapter "job-copilot-backend/internal/adapter/http"
	"job-copilot-backend/internal/adapter/repository"
	analysisapp "job-copilot-backend/internal/application/analysis"
	interviewapp "job-copilot-backend/internal/application/interview"
	"job-copilot-backend/internal/config"
	"job-copilot-backend/internal/port"
)

func main() {
	appConfig := config.Load()

	var authProvider port.AuthProvider
	if configuredAuth, err := auth.NewSupabaseAuthAdapter(
		appConfig.SupabaseURL,
		appConfig.SupabaseAnonKey,
		nil,
	); err == nil {
		authProvider = configuredAuth
	}

	var aiClient port.AIClient = ai.NewPlaceholderAdapter()
	if configuredAI, err := ai.NewOpenAICompatibleAdapter(
		appConfig.AIBaseURL,
		appConfig.AIAPIKey,
		appConfig.AIModel,
		&http.Client{Timeout: appConfig.AITimeout},
	); err == nil {
		aiClient = configuredAI
	}

	var analysisRepository port.AnalysisRepository = repository.NewPlaceholderAnalysisRepository()
	if configuredRepository, err := repository.NewSupabaseAnalysisRepository(
		appConfig.SupabaseURL,
		appConfig.SupabaseAnonKey,
		nil,
	); err == nil {
		analysisRepository = configuredRepository
	}
	var interviewRepository port.InterviewRepository = repository.NewPlaceholderInterviewRepository()
	if configuredRepository, err := repository.NewSupabaseInterviewRepository(
		appConfig.SupabaseURL,
		appConfig.SupabaseAnonKey,
		nil,
	); err == nil {
		interviewRepository = configuredRepository
	}

	analyzeJDService, err := analysisapp.NewAnalyzeJDService(aiClient, analysisRepository)
	if err != nil {
		log.Fatalf("初始化 JD 分析服务失败：%v", err)
	}
	analysisHistoryService, err := analysisapp.NewHistoryService(analysisRepository)
	if err != nil {
		log.Fatalf("初始化 JD 历史服务失败：%v", err)
	}
	startInterviewService, err := interviewapp.NewStartInterviewService(aiClient, interviewRepository)
	if err != nil {
		log.Fatalf("初始化模拟面试服务失败：%v", err)
	}
	listInterviewOptionsService, err := interviewapp.NewListInterviewOptionsService(interviewRepository)
	if err != nil {
		log.Fatalf("初始化面试岗位选项服务失败：%v", err)
	}
	interviewTurnService, err := interviewapp.NewTurnService(aiClient, interviewRepository)
	if err != nil {
		log.Fatalf("初始化面试单轮服务失败：%v", err)
	}
	interviewSessionService, err := interviewapp.NewSessionService(interviewRepository)
	if err != nil {
		log.Fatalf("初始化面试会话服务失败：%v", err)
	}
	interviewReportService, err := interviewapp.NewReportService(aiClient, interviewRepository)
	if err != nil {
		log.Fatalf("初始化面试报告服务失败：%v", err)
	}
	engine := httpadapter.NewRouter(appConfig.FrontendOrigin, httpadapter.Dependencies{
		AuthProvider:                authProvider,
		AnalyzeJDService:            analyzeJDService,
		AnalysisHistoryService:      analysisHistoryService,
		StartInterviewService:       startInterviewService,
		ListInterviewOptionsService: listInterviewOptionsService,
		InterviewTurnService:        interviewTurnService,
		InterviewSessionService:     interviewSessionService,
		InterviewReportService:      interviewReportService,
	})
	address := ":" + appConfig.AppPort

	log.Printf(
		"Job Copilot API 正在启动：环境=%s，地址=http://localhost%s",
		appConfig.AppEnv,
		address,
	)

	if err := engine.Run(address); err != nil {
		log.Fatalf("启动服务失败：%v", err)
	}
}
