package analysis

import (
	"context"
	"errors"
	"strings"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	"job-copilot-backend/internal/port"
)

var (
	ErrMissingDependency = errors.New("analyze JD service dependency is missing")
	ErrMissingUserID     = errors.New("authenticated user ID is required")
)

type AnalyzeJDService struct {
	aiClient   port.AIClient
	repository port.AnalysisRepository
}

type AnalyzeJDOutput struct {
	AnalysisID string
	Result     analysisdomain.AnalysisResult
}

func NewAnalyzeJDService(
	aiClient port.AIClient,
	repository port.AnalysisRepository,
) (*AnalyzeJDService, error) {
	if aiClient == nil || repository == nil {
		return nil, ErrMissingDependency
	}

	return &AnalyzeJDService{
		aiClient:   aiClient,
		repository: repository,
	}, nil
}

// Execute 只编排用例顺序；AI 调用和数据保存由 Port 的具体 Adapter 执行。
func (service *AnalyzeJDService) Execute(
	ctx context.Context,
	userID string,
	jdContent string,
) (AnalyzeJDOutput, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return AnalyzeJDOutput{}, ErrMissingUserID
	}

	description, err := analysisdomain.NewJobDescription(jdContent)
	if err != nil {
		return AnalyzeJDOutput{}, err
	}

	result, err := service.aiClient.AnalyzeJD(ctx, description)
	if err != nil {
		return AnalyzeJDOutput{}, err
	}

	analysisID, err := service.repository.Save(ctx, userID, result)
	if err != nil {
		return AnalyzeJDOutput{}, err
	}

	return AnalyzeJDOutput{
		AnalysisID: analysisID,
		Result:     result,
	}, nil
}
