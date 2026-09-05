package port

import (
	"context"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
)

// AIClient 描述 Application 需要的 JD 分析能力，不暴露具体模型 SDK。
type AIClient interface {
	AnalyzeJD(ctx context.Context, description analysisdomain.JobDescription) (analysisdomain.AnalysisResult, error)
}
