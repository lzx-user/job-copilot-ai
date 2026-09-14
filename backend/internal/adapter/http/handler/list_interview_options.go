package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	interviewapp "job-copilot-backend/internal/application/interview"
	"job-copilot-backend/pkg/response"
)

type ListInterviewOptionsHandler struct {
	service *interviewapp.ListInterviewOptionsService
}

func NewListInterviewOptionsHandler(service *interviewapp.ListInterviewOptionsService) *ListInterviewOptionsHandler {
	return &ListInterviewOptionsHandler{service: service}
}

func (handler *ListInterviewOptionsHandler) Handle(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}

	options, err := handler.service.Execute(ctx.Request.Context(), userID)
	if err != nil {
		handleStartInterviewError(ctx, err)
		return
	}

	items := make([]gin.H, 0, len(options))
	for _, option := range options {
		items = append(items, gin.H{
			"analysisId":  option.ID(),
			"companyName": option.CompanyName(),
			"jobTitle":    option.JobTitle(),
			"matchScore":  option.MatchScore(),
			"createdAt":   option.CreatedAt(),
		})
	}
	response.Success(ctx, gin.H{"items": items})
}
