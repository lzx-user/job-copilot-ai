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

const maxInterviewReportRequestBytes = 8 << 10

type InterviewReportHandler struct{ service *interviewapp.ReportService }

func NewInterviewReportHandler(service *interviewapp.ReportService) *InterviewReportHandler {
	return &InterviewReportHandler{service: service}
}

func (handler *InterviewReportHandler) Handle(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}
	var request struct {
		SessionID string `json:"sessionId"`
	}
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxInterviewReportRequestBytes)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "请求体必须是有效的 JSON")
		return
	}
	output, err := handler.service.Execute(ctx.Request.Context(), userID, request.SessionID)
	if err != nil {
		handleInterviewReportError(ctx, err)
		return
	}
	data := interviewReportData(output.Report)
	data["sessionId"] = output.SessionID
	data["status"] = output.Status
	response.Success(ctx, data)
}

func interviewReportData(report interviewdomain.InterviewReport) gin.H {
	return gin.H{
		"overallScore": report.OverallScore(), "technicalScore": report.TechnicalScore(),
		"expressionScore": report.ExpressionScore(), "projectDepthScore": report.ProjectDepthScore(),
		"strengths": report.Strengths(), "weaknesses": report.Weaknesses(),
		"recommendedTopics": report.RecommendedTopics(), "answerTips": report.AnswerTips(),
		"summary": report.Summary(),
	}
}

func handleInterviewReportError(ctx *gin.Context, err error) {
	log.Printf("interview report failed: %v", err)
	switch {
	case errors.Is(err, interviewapp.ErrInvalidInterviewReportCommand):
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "会话 ID 无效")
	case errors.Is(err, interviewapp.ErrInterviewReportUnavailable), errors.Is(err, port.ErrRepositoryConflict):
		response.Error(ctx, http.StatusConflict, "INTERVIEW_REPORT_CONFLICT", "请先完成全部五轮问答，或刷新查看已生成的报告")
	case errors.Is(err, port.ErrRepositoryNotFound):
		response.Error(ctx, http.StatusNotFound, "NOT_FOUND", "未找到该面试会话")
	case errors.Is(err, port.ErrAIUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "AI_NOT_CONFIGURED", "AI 服务未配置")
	case errors.Is(err, context.DeadlineExceeded):
		response.Error(ctx, http.StatusGatewayTimeout, "AI_TIMEOUT", "最终报告生成超时，请稍后重试")
	case errors.Is(err, port.ErrAIInvalidResponse):
		response.Error(ctx, http.StatusBadGateway, "AI_INVALID_RESPONSE", "AI 返回的最终报告无法使用")
	case errors.Is(err, port.ErrAIUpstream):
		response.Error(ctx, http.StatusBadGateway, "AI_UPSTREAM_ERROR", "AI 服务暂时不可用")
	case errors.Is(err, port.ErrRepositoryUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
	default:
		response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "保存最终报告失败，请稍后重试")
	}
}
