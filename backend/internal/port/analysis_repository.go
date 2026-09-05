package port

import (
	"context"

	"job-copilot-backend/internal/domain"
)

// AnalysisRepository 描述核心业务需要的分析结果持久化能力。
type AnalysisRepository interface {
	Save(ctx context.Context, userID string, result domain.AnalysisResult) (analysisID string, err error)
	FindByID(ctx context.Context, userID string, analysisID string) (domain.AnalysisResult, error)
}
