package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	interviewapp "job-copilot-backend/internal/application/interview"
	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
	"job-copilot-backend/pkg/response"
)

type InterviewHistoryHandler struct{ service *interviewapp.HistoryService }

func NewInterviewHistoryHandler(service *interviewapp.HistoryService) *InterviewHistoryHandler {
	return &InterviewHistoryHandler{service: service}
}

func (handler *InterviewHistoryHandler) List(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}
	records, err := handler.service.List(ctx.Request.Context(), userID)
	if err != nil {
		handleInterviewHistoryError(ctx, err)
		return
	}
	items := make([]gin.H, 0, len(records))
	for _, record := range records {
		items = append(items, interviewHistoryResponse(record))
	}
	response.Success(ctx, gin.H{"items": items})
}

func (handler *InterviewHistoryHandler) Delete(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}
	if err := handler.service.Delete(ctx.Request.Context(), userID, ctx.Param("id")); err != nil {
		handleInterviewHistoryError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"sessionId": ctx.Param("id")})
}

func interviewHistoryResponse(record interviewdomain.InterviewHistoryRecord) gin.H {
	return gin.H{
		"sessionId": record.ID(), "companyName": record.CompanyName(), "jobTitle": record.JobTitle(),
		"status": record.Status(), "currentRound": record.CurrentRound(), "maxRounds": record.MaxRounds(),
		"overallScore": record.OverallScore(), "createdAt": record.CreatedAt(),
		"updatedAt": record.UpdatedAt(), "completedAt": record.CompletedAt(),
	}
}

func handleInterviewHistoryError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, interviewapp.ErrInvalidInterviewSessionQuery):
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "会话 ID 无效")
	case errors.Is(err, port.ErrRepositoryNotFound):
		response.Error(ctx, http.StatusNotFound, "NOT_FOUND", "未找到该面试会话")
	case errors.Is(err, port.ErrRepositoryUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
	default:
		response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "读取面试历史失败，请稍后重试")
	}
}
