package main

import (
	"log"

	httpadapter "job-copilot-backend/internal/adapter/http"
	"job-copilot-backend/internal/config"
)

func main() {
	appConfig := config.Load()
	engine := httpadapter.NewRouter(appConfig.FrontendOrigin)
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
