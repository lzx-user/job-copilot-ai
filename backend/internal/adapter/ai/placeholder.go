package ai

import (
	"context"
	"errors"

	"job-copilot-backend/internal/domain"
	"job-copilot-backend/internal/port"
)

var ErrClientNotConfigured = errors.New("AI client is not configured")

// PlaceholderAdapter 只用于保持当前架构可组装；它不会生成伪造的分析结果。
type PlaceholderAdapter struct{}

var _ port.AIClient = (*PlaceholderAdapter)(nil)

func NewPlaceholderAdapter() *PlaceholderAdapter {
	return &PlaceholderAdapter{}
}

func (adapter *PlaceholderAdapter) AnalyzeJD(
	ctx context.Context,
	_ domain.JobDescription,
) (domain.AnalysisResult, error) {
	if err := ctx.Err(); err != nil {
		return domain.AnalysisResult{}, err
	}

	return domain.AnalysisResult{}, ErrClientNotConfigured
}
