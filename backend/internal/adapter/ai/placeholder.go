package ai

import (
	"context"
	"errors"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
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
	_ analysisdomain.JobDescription,
) (analysisdomain.AnalysisResult, error) {
	if err := ctx.Err(); err != nil {
		return analysisdomain.AnalysisResult{}, err
	}

	return analysisdomain.AnalysisResult{}, ErrClientNotConfigured
}
