package handler

import (
	"github.com/gin-gonic/gin"

	"job-copilot-backend/pkg/response"
)

func Health(ctx *gin.Context) {
	response.Success(ctx, gin.H{
		"status": "ok",
	})
}
