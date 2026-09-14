package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	analysisapp "job-copilot-backend/internal/application/analysis"
	analysisdomain "job-copilot-backend/internal/domain/analysis"
	"job-copilot-backend/internal/port"
	"job-copilot-backend/pkg/response"
)

type JDAnalysisHistoryHandler struct {
	service *analysisapp.HistoryService
}

func NewJDAnalysisHistoryHandler(service *analysisapp.HistoryService) *JDAnalysisHistoryHandler {
	return &JDAnalysisHistoryHandler{service: service}
}

func (handler *JDAnalysisHistoryHandler) List(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}
	records, err := handler.service.List(ctx.Request.Context(), userID)
	if err != nil {
		handleAnalysisHistoryError(ctx, err)
		return
	}
	items := make([]gin.H, 0, len(records))
	for _, record := range records {
		items = append(items, analysisRecordResponse(record))
	}
	response.Success(ctx, gin.H{"items": items})
}

func (handler *JDAnalysisHistoryHandler) Get(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}
	record, err := handler.service.Get(ctx.Request.Context(), userID, ctx.Param("id"))
	if err != nil {
		handleAnalysisHistoryError(ctx, err)
		return
	}
	response.Success(ctx, analysisRecordResponse(record))
}

func analysisRecordResponse(record analysisdomain.AnalysisRecord) gin.H {
	result := record.Result()
	return gin.H{
		"analysisId": record.ID(), "companyName": record.CompanyName(),
		"jobTitle": record.JobTitle(), "createdAt": record.CreatedAt(),
		"matchScore": result.MatchScore(), "jobSummary": result.JobSummary(),
		"coreRequirements": result.CoreRequirements(), "matchedSkills": result.MatchedSkills(),
		"missingSkills": result.MissingSkills(), "resumeSuggestions": result.ResumeSuggestions(),
		"preparationTopics": result.PreparationTopics(), "greetingMessage": result.GreetingMessage(),
	}
}

func handleAnalysisHistoryError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, analysisapp.ErrInvalidAnalysisHistoryQuery):
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "分析记录 ID 无效")
	case errors.Is(err, port.ErrRepositoryNotFound):
		response.Error(ctx, http.StatusNotFound, "NOT_FOUND", "未找到该分析记录")
	case errors.Is(err, port.ErrRepositoryUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
	default:
		response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "读取分析记录失败，请稍后重试")
	}
}
