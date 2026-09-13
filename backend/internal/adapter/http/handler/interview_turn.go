package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	interviewapp "job-copilot-backend/internal/application/interview"
	"job-copilot-backend/internal/port"
	"job-copilot-backend/pkg/response"
)

const maxInterviewTurnRequestBytes = 16 << 10

type InterviewTurnHandler struct{ service *interviewapp.TurnService }

func NewInterviewTurnHandler(service *interviewapp.TurnService) *InterviewTurnHandler {
	return &InterviewTurnHandler{service: service}
}

func (handler *InterviewTurnHandler) Handle(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}
	var request struct {
		SessionID string `json:"sessionId"`
		Answer    string `json:"answer"`
	}
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxInterviewTurnRequestBytes)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "请求体必须是有效的 JSON")
		return
	}
	output, err := handler.service.Execute(ctx.Request.Context(), userID, interviewapp.TurnCommand{
		SessionID: request.SessionID, Answer: request.Answer,
	})
	if err != nil {
		handleInterviewTurnError(ctx, err)
		return
	}
	feedback := output.Result.Feedback()
	var nextQuestion any
	if output.Result.HasNextQuestion() {
		nextQuestion = output.Result.NextQuestion()
	}
	response.Success(ctx, gin.H{
		"sessionId": output.SessionID, "currentRound": output.CurrentRound, "maxRounds": output.MaxRounds,
		"score": feedback.Score(), "feedback": feedback.Feedback(), "strengths": feedback.Strengths(),
		"improvements": feedback.Improvements(), "nextQuestion": nextQuestion,
	})
}

func handleInterviewTurnError(ctx *gin.Context, err error) {
	log.Printf("interview turn failed: %v", err)
	switch {
	case errors.Is(err, interviewapp.ErrInvalidInterviewTurnCommand):
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "请检查会话 ID 和回答内容")
	case errors.Is(err, interviewapp.ErrInterviewTurnUnavailable), errors.Is(err, port.ErrRepositoryConflict):
		response.Error(ctx, http.StatusConflict, "INTERVIEW_CONFLICT", "当前轮次已提交或暂不能继续")
	case errors.Is(err, port.ErrRepositoryNotFound):
		response.Error(ctx, http.StatusNotFound, "NOT_FOUND", "未找到该面试会话")
	case errors.Is(err, port.ErrAIUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "AI_NOT_CONFIGURED", "AI 服务未配置")
	case errors.Is(err, context.DeadlineExceeded):
		response.Error(ctx, http.StatusGatewayTimeout, "AI_TIMEOUT", "AI 评价超时，请稍后重试")
	case errors.Is(err, port.ErrAIInvalidResponse):
		response.Error(ctx, http.StatusBadGateway, "AI_INVALID_RESPONSE", "AI 返回的评价无法使用")
	case errors.Is(err, port.ErrAIUpstream):
		response.Error(ctx, http.StatusBadGateway, "AI_UPSTREAM_ERROR", "AI 服务暂时不可用")
	case errors.Is(err, port.ErrRepositoryUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
	default:
		response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "保存面试回答失败，请稍后重试")
	}
}
