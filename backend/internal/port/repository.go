package port

import (
	"context"
	"errors"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
)

var (
	ErrRepositoryUnavailable = errors.New("repository unavailable")
	ErrRepositoryOperation   = errors.New("repository operation failed")
)

// AnalysisRepository 描述核心业务需要的分析结果持久化能力。
type AnalysisRepository interface {
	Save(
		ctx context.Context,
		userID string,
		request analysisdomain.AnalysisRequest,
		result analysisdomain.AnalysisResult,
	) (analysisID string, err error)
	FindByID(ctx context.Context, userID string, analysisID string) (analysisdomain.AnalysisResult, error)
}
