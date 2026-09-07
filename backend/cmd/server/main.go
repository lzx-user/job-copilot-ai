package main

import (
	"log"
	"net/http"

	"job-copilot-backend/internal/adapter/ai"
	"job-copilot-backend/internal/adapter/auth"
	httpadapter "job-copilot-backend/internal/adapter/http"
	analysisapp "job-copilot-backend/internal/application/analysis"
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

	analyzeJDService, err := analysisapp.NewAnalyzeJDService(aiClient)
	if err != nil {
		log.Fatalf("初始化 JD 分析服务失败：%v", err)
	}
	engine := httpadapter.NewRouter(appConfig.FrontendOrigin, httpadapter.Dependencies{
		AuthProvider:     authProvider,
		AnalyzeJDService: analyzeJDService,
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
