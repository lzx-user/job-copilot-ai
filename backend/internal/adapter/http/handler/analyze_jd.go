package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"job-copilot-backend/internal/adapter/http/middleware"
	analysisapp "job-copilot-backend/internal/application/analysis"
	analysisdomain "job-copilot-backend/internal/domain/analysis"
	"job-copilot-backend/internal/port"
	"job-copilot-backend/pkg/response"
)

const maxAnalyzeJDRequestBytes = 64 << 10

type AnalyzeJDHandler struct {
	service *analysisapp.AnalyzeJDService
}

type analyzeJDRequest struct {
	CompanyName   string   `json:"companyName"`
	JobTitle      string   `json:"jobTitle"`
	JDContent     string   `json:"jdContent"`
	ResumeSummary string   `json:"resumeSummary"`
	Skills        []string `json:"skills"`
}

func NewAnalyzeJDHandler(service *analysisapp.AnalyzeJDService) *AnalyzeJDHandler {
	return &AnalyzeJDHandler{service: service}
}

func (handler *AnalyzeJDHandler) Handle(ctx *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
		return
	}

	var request analyzeJDRequest
	// 在 JSON 解码前限制请求体，避免超大输入消耗服务器内存和 AI 成本。
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxAnalyzeJDRequestBytes)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "请求体必须是有效的 JSON")
		return
	}

	output, err := handler.service.Execute(ctx.Request.Context(), userID, analysisapp.AnalyzeJDCommand{
		CompanyName:   request.CompanyName,
		JobTitle:      request.JobTitle,
		JDContent:     request.JDContent,
		ResumeSummary: request.ResumeSummary,
		Skills:        request.Skills,
	})
	if err != nil {
		handleAnalyzeJDError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"analysisId":        output.AnalysisID,
		"matchScore":        output.Result.MatchScore(),
		"jobSummary":        output.Result.JobSummary(),
		"coreRequirements":  output.Result.CoreRequirements(),
		"matchedSkills":     output.Result.MatchedSkills(),
		"missingSkills":     output.Result.MissingSkills(),
		"resumeSuggestions": output.Result.ResumeSuggestions(),
		"preparationTopics": output.Result.PreparationTopics(),
		"greetingMessage":   output.Result.GreetingMessage(),
	})
}

func handleAnalyzeJDError(ctx *gin.Context, err error) {
	// 只记录错误链，不记录用户标识、JD、简历、Token 或上游响应体。
	log.Printf("JD analysis request failed: %v", err)

	switch {
	case errors.Is(err, analysisdomain.ErrInvalidAnalysisRequest), errors.Is(err, analysisdomain.ErrEmptyJobDescription):
		response.Error(ctx, http.StatusBadRequest, "INVALID_ARGUMENT", "请检查公司、岗位、JD 和个人经历信息")
	case errors.Is(err, port.ErrAIUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "AI_NOT_CONFIGURED", "AI 服务未配置")
	case errors.Is(err, context.DeadlineExceeded):
		response.Error(ctx, http.StatusGatewayTimeout, "AI_TIMEOUT", "AI 分析超时，请稍后重试")
	case errors.Is(err, port.ErrAIInvalidResponse):
		response.Error(ctx, http.StatusBadGateway, "AI_INVALID_RESPONSE", "AI 返回的结果无法使用")
	case errors.Is(err, port.ErrAIUpstream):
		response.Error(ctx, http.StatusBadGateway, "AI_UPSTREAM_ERROR", "AI 服务暂时不可用")
	case errors.Is(err, port.ErrRepositoryUnavailable):
		response.Error(ctx, http.StatusServiceUnavailable, "DATABASE_NOT_CONFIGURED", "数据存储服务未配置")
	case errors.Is(err, port.ErrRepositoryOperation):
		response.Error(ctx, http.StatusBadGateway, "DATABASE_ERROR", "分析结果保存失败，请稍后重试")
	default:
		response.Error(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "分析失败，请稍后重试")
	}
}
