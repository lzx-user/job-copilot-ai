package port

import (
	"context"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
)

// AnalysisRepository 描述核心业务需要的分析结果持久化能力。
type AnalysisRepository interface {
	Save(ctx context.Context, userID string, result analysisdomain.AnalysisResult) (analysisID string, err error)
	FindByID(ctx context.Context, userID string, analysisID string) (analysisdomain.AnalysisResult, error)
}
