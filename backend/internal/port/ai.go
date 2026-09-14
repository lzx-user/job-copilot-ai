package port

import (
	"context"
	"errors"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	interviewdomain "job-copilot-backend/internal/domain/interview"
)

var (
	ErrAIUnavailable     = errors.New("AI service unavailable")
	ErrAIUpstream        = errors.New("AI upstream request failed")
	ErrAIInvalidResponse = errors.New("AI returned an invalid response")
)

// AIClient 描述 Application 当前需要的 AI 能力，不暴露具体模型 SDK。
type AIClient interface {
	AnalyzeJD(ctx context.Context, request analysisdomain.AnalysisRequest) (analysisdomain.AnalysisResult, error)
	GenerateFirstInterviewQuestion(ctx context.Context, input interviewdomain.InterviewContext) (string, error)
	EvaluateInterviewAnswer(ctx context.Context, input interviewdomain.InterviewTurnPrompt) (interviewdomain.InterviewTurnResult, error)
	GenerateInterviewReport(ctx context.Context, input interviewdomain.InterviewReportContext) (interviewdomain.InterviewReport, error)
}
