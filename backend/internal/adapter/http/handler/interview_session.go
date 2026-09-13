package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	interviewapp "job-copilot-backend/internal/application/interview"
	"job-copilot-backend/internal/port"
	"job-copilot-backend/pkg/response"
)

type InterviewSessionHandler struct{ service *interviewapp.SessionService }

func NewInterviewSessionHandler(service *interviewapp.SessionService) *InterviewSessionHandler {
	return &InterviewSessionHandler{service: service}
}

func (handler *InterviewSessionHandler) Get(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}
	detail, err := handler.service.Get(ctx.Request.Context(), userID, ctx.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, interviewapp.ErrInvalidInterviewSessionQuery):
			response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "会话 ID 无效")
		case errors.Is(err, port.ErrRepositoryNotFound):
			response.Error(ctx, http.StatusNotFound, "NOT_FOUND", "未找到该面试会话")
		case errors.Is(err, port.ErrRepositoryUnavailable):
			response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
		default:
			response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "读取面试会话失败，请稍后重试")
		}
		return
	}
	messages := make([]gin.H, 0, len(detail.Messages))
	for _, message := range detail.Messages {
		item := gin.H{
			"id": message.ID, "role": message.Role, "round": message.Round,
			"content": message.Content, "createdAt": message.CreatedAt,
		}
		if message.Feedback != nil {
			item["score"] = message.Feedback.Score()
			item["feedback"] = message.Feedback.Feedback()
			item["strengths"] = message.Feedback.Strengths()
			item["improvements"] = message.Feedback.Improvements()
		}
		messages = append(messages, item)
	}
	response.Success(ctx, gin.H{
		"sessionId": detail.Session.ID(), "status": detail.Session.Status(),
		"currentRound": detail.Session.CurrentRound(), "maxRounds": detail.Session.MaxRounds(),
		"companyName": detail.CompanyName, "jobTitle": detail.JobTitle, "messages": messages,
	})
}
