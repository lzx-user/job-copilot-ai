package handler

import (
	"github.com/gin-gonic/gin"

	"job-copilot-server/pkg/response"
)

func Health(ctx *gin.Context) {
	response.Success(ctx, gin.H{
		"status": "ok",
	})
}
