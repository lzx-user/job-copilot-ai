package ai

import (
	"context"
	"errors"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

var ErrClientNotConfigured = errors.New("AI client is not configured")

// PlaceholderAdapter 只用于保持当前架构可组装；它不会生成伪造的分析结果。
type PlaceholderAdapter struct{}

var _ port.AIClient = (*PlaceholderAdapter)(nil)

func NewPlaceholderAdapter() *PlaceholderAdapter {
	return &PlaceholderAdapter{}
}

func (adapter *PlaceholderAdapter) GenerateFirstInterviewQuestion(
	ctx context.Context,
	_ interviewdomain.InterviewContext,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	return "", errors.Join(port.ErrAIUnavailable, ErrClientNotConfigured)
}

func (adapter *PlaceholderAdapter) AnalyzeJD(
	ctx context.Context,
	_ analysisdomain.AnalysisRequest,
) (analysisdomain.AnalysisResult, error) {
	if err := ctx.Err(); err != nil {
		return analysisdomain.AnalysisResult{}, err
	}

	return analysisdomain.AnalysisResult{}, errors.Join(port.ErrAIUnavailable, ErrClientNotConfigured)
}
