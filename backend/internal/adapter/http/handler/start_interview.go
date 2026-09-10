package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	interviewapp "job-copilot-backend/internal/application/interview"
	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
	"job-copilot-backend/pkg/response"
)

const maxStartInterviewRequestBytes = 4 << 10

type StartInterviewHandler struct {
	service *interviewapp.StartInterviewService
}

type startInterviewRequest struct {
	AnalysisID string `json:"analysisId"`
}

func NewStartInterviewHandler(service *interviewapp.StartInterviewService) *StartInterviewHandler {
	return &StartInterviewHandler{service: service}
}

func (handler *StartInterviewHandler) Handle(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}

	var request startInterviewRequest
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxStartInterviewRequestBytes)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "请求体必须是有效的 JSON")
		return
	}

	output, err := handler.service.Execute(ctx.Request.Context(), userID, interviewapp.StartInterviewCommand{
		AnalysisID: request.AnalysisID,
	})
	if err != nil {
		handleStartInterviewError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"sessionId":    output.SessionID,
		"status":       output.Status,
		"currentRound": output.CurrentRound,
		"maxRounds":    output.MaxRounds,
		"question":     output.Question,
	})
}

func handleStartInterviewError(ctx *gin.Context, err error) {
	// 不记录用户标识、JD、简历、Token、Prompt 或生成的问题正文。
	log.Printf("start interview request failed: %v", err)

	switch {
	case errors.Is(err, interviewapp.ErrInvalidStartInterviewCommand):
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "请选择有效的 JD 分析记录")
	case errors.Is(err, port.ErrRepositoryNotFound):
		response.Error(ctx, http.StatusNotFound, "JD_ANALYSIS_NOT_FOUND", "未找到可用于面试的 JD 分析记录")
	case errors.Is(err, port.ErrAIUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "AI_NOT_CONFIGURED", "AI 服务未配置")
	case errors.Is(err, context.DeadlineExceeded):
		response.Error(ctx, http.StatusGatewayTimeout, "AI_TIMEOUT", "面试题生成超时，请稍后重试")
	case errors.Is(err, port.ErrAIInvalidResponse):
		response.Error(ctx, http.StatusBadGateway, "AI_INVALID_RESPONSE", "AI 返回的面试题无法使用")
	case errors.Is(err, port.ErrAIUpstream):
		response.Error(ctx, http.StatusBadGateway, "AI_UPSTREAM_ERROR", "AI 服务暂时不可用")
	case errors.Is(err, port.ErrRepositoryUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
	case errors.Is(err, port.ErrRepositoryOperation),
		errors.Is(err, interviewdomain.ErrInvalidInterviewContext),
		errors.Is(err, interviewdomain.ErrInvalidInterviewSession):
		response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "面试会话保存失败，请稍后重试")
	default:
		response.Error(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "面试启动失败，请稍后重试")
	}
}
