package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-server/pkg/response"
)

func NotFound(ctx *gin.Context) {
	response.Error(
		ctx,
		http.StatusNotFound,
		"NOT_FOUND",
		"未找到请求的接口",
	)
}
