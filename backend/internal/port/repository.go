package port

import (
	"context"
	"errors"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	careerdomain "job-copilot-backend/internal/domain/career"
	dashboarddomain "job-copilot-backend/internal/domain/dashboard"
	interviewdomain "job-copilot-backend/internal/domain/interview"
)

type CareerRepository interface {
	GetWorkspace(ctx context.Context, userID string) (careerdomain.Workspace, error)
	CreateResume(ctx context.Context, userID string, input careerdomain.CreateResumeInput) (resumeID, versionID string, err error)
	CreateResumeVersion(ctx context.Context, userID string, input careerdomain.CreateResumeVersionInput) (string, error)
	CreateJob(ctx context.Context, userID string, input careerdomain.CreateJobInput) (jobID, jdVersionID string, err error)
	CreateJDVersion(ctx context.Context, userID string, input careerdomain.CreateJDVersionInput) (string, error)
	CreateApplication(ctx context.Context, userID string, input careerdomain.CreateApplicationInput) (string, error)
	AddApplicationEvent(ctx context.Context, userID string, input careerdomain.CreateEventInput) (string, error)
	CreateRealInterview(ctx context.Context, userID string, input careerdomain.CreateInterviewInput) (string, error)
	UpdateRealInterview(ctx context.Context, userID string, input careerdomain.UpdateInterviewInput) error
	SaveRetrospective(ctx context.Context, userID string, input careerdomain.SaveRetrospectiveInput) error
	SaveOffer(ctx context.Context, userID string, input careerdomain.SaveOfferInput) error
}

type DashboardRepository interface {
	GetSummary(ctx context.Context, userID string) (dashboarddomain.Summary, error)
}

var (
	ErrRepositoryUnavailable = errors.New("repository unavailable")
	ErrRepositoryOperation   = errors.New("repository operation failed")
	ErrRepositoryNotFound    = errors.New("repository record not found")
	ErrRepositoryConflict    = errors.New("repository conflict")
	ErrRepositorySchema      = errors.New("repository schema is not initialized")
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
	List(ctx context.Context, userID string, limit int) ([]analysisdomain.AnalysisRecord, error)
	FindRecordByID(ctx context.Context, userID string, analysisID string) (analysisdomain.AnalysisRecord, error)
	Delete(ctx context.Context, userID string, analysisID string) error
}

// InterviewRepository 只暴露五轮面试启动、恢复与持久化所需的能力。
type InterviewRepository interface {
	ListOptions(ctx context.Context, userID string) ([]interviewdomain.InterviewOption, error)
	FindContext(ctx context.Context, userID string, analysisID string) (interviewdomain.InterviewContext, error)
	CreatePendingSession(ctx context.Context, session interviewdomain.InterviewSession) (sessionID string, err error)
	StartSession(
		ctx context.Context,
		session interviewdomain.InterviewSession,
		question interviewdomain.InterviewMessage,
	) error
	FindTurnContext(ctx context.Context, userID string, sessionID string) (interviewdomain.InterviewTurnContext, error)
	FindReportContext(ctx context.Context, userID string, sessionID string) (interviewdomain.InterviewReportContext, error)
	FindSessionDetail(ctx context.Context, userID string, sessionID string) (interviewdomain.SessionDetail, error)
	ListSessions(ctx context.Context, userID string, limit int) ([]interviewdomain.InterviewHistoryRecord, error)
	DeleteSession(ctx context.Context, userID string, sessionID string) error
	SaveTurn(
		ctx context.Context,
		session interviewdomain.InterviewSession,
		answer interviewdomain.InterviewMessage,
		result interviewdomain.InterviewTurnResult,
	) error
	SaveReport(ctx context.Context, session interviewdomain.InterviewSession, report interviewdomain.InterviewReport) error
}
