package port

import (
	"context"
	"errors"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
)

var (
	ErrAIUnavailable     = errors.New("AI service unavailable")
	ErrAIUpstream        = errors.New("AI upstream request failed")
	ErrAIInvalidResponse = errors.New("AI returned an invalid response")
)

// AIClient 描述 Application 需要的 JD 分析能力，不暴露具体模型 SDK。
type AIClient interface {
	AnalyzeJD(ctx context.Context, request analysisdomain.AnalysisRequest) (analysisdomain.AnalysisResult, error)
}
