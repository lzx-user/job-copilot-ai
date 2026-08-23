package main

import (
	"log"

	httpadapter "job-copilot-server/internal/adapter/http"
	"job-copilot-server/internal/config"
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
