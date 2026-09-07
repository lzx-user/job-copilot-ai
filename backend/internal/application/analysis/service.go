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
	aiClient port.AIClient
}

type AnalyzeJDOutput struct {
	Result analysisdomain.AnalysisResult
}

type AnalyzeJDCommand struct {
	CompanyName   string
	JobTitle      string
	JDContent     string
	ResumeSummary string
	Skills        []string
}

func NewAnalyzeJDService(aiClient port.AIClient) (*AnalyzeJDService, error) {
	if aiClient == nil {
		return nil, ErrMissingDependency
	}

	return &AnalyzeJDService{aiClient: aiClient}, nil
}

// Execute 只编排当前用例顺序；AI 调用由 Port 的具体 Adapter 执行，持久化在后续阶段接入。
func (service *AnalyzeJDService) Execute(
	ctx context.Context,
	userID string,
	command AnalyzeJDCommand,
) (AnalyzeJDOutput, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return AnalyzeJDOutput{}, ErrMissingUserID
	}

	request, err := analysisdomain.NewAnalysisRequest(
		command.CompanyName,
		command.JobTitle,
		command.JDContent,
		command.ResumeSummary,
		command.Skills,
	)
	if err != nil {
		return AnalyzeJDOutput{}, err
	}

	result, err := service.aiClient.AnalyzeJD(ctx, request)
	if err != nil {
		return AnalyzeJDOutput{}, err
	}

	return AnalyzeJDOutput{Result: result}, nil
}
