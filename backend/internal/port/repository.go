package port

import (
	"context"
	"errors"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	interviewdomain "job-copilot-backend/internal/domain/interview"
)

var (
	ErrRepositoryUnavailable = errors.New("repository unavailable")
	ErrRepositoryOperation   = errors.New("repository operation failed")
	ErrRepositoryNotFound    = errors.New("repository record not found")
	ErrRepositoryConflict    = errors.New("repository conflict")
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

// InterviewRepository 只暴露启动面试所需的读取与持久化能力。
type InterviewRepository interface {
	ListOptions(ctx context.Context, userID string) ([]interviewdomain.InterviewOption, error)
	FindContext(ctx context.Context, userID string, analysisID string) (interviewdomain.InterviewContext, error)
	CreatePendingSession(ctx context.Context, userID string, analysisID string) (sessionID string, err error)
	StartSession(
		ctx context.Context,
		session interviewdomain.InterviewSession,
		question interviewdomain.InterviewMessage,
	) error
	FindTurnContext(ctx context.Context, userID string, sessionID string) (interviewdomain.InterviewTurnContext, error)
	FindSessionDetail(ctx context.Context, userID string, sessionID string) (interviewdomain.SessionDetail, error)
	SaveTurn(
		ctx context.Context,
		session interviewdomain.InterviewSession,
		answer interviewdomain.InterviewMessage,
		result interviewdomain.InterviewTurnResult,
	) error
}
