package interview

import (
	"context"
	"errors"
	"strings"

	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

var (
	ErrInvalidInterviewReportCommand = errors.New("invalid interview report command")
	ErrInterviewReportUnavailable    = errors.New("interview report is unavailable")
)

type ReportService struct {
	aiClient   port.AIClient
	repository port.InterviewRepository
}

type ReportOutput struct {
	SessionID string
	Status    interviewdomain.InterviewStatus
	Report    interviewdomain.InterviewReport
}

func NewReportService(aiClient port.AIClient, repository port.InterviewRepository) (*ReportService, error) {
	if aiClient == nil || repository == nil {
		return nil, ErrMissingDependency
	}
	return &ReportService{aiClient: aiClient, repository: repository}, nil
}

func (service *ReportService) Execute(ctx context.Context, userID, sessionID string) (ReportOutput, error) {
	userID = strings.TrimSpace(userID)
	sessionID = strings.TrimSpace(sessionID)
	if userID == "" || !isUUID(sessionID) {
		return ReportOutput{}, ErrInvalidInterviewReportCommand
	}

	reportContext, err := service.repository.FindReportContext(ctx, userID, sessionID)
	if err != nil {
		return ReportOutput{}, err
	}
	// 已完成时直接返回持久化报告，避免网络重试触发第二次模型调用。
	if reportContext.Session.Status() == interviewdomain.InterviewStatusCompleted && reportContext.Report != nil {
		return ReportOutput{SessionID: sessionID, Status: reportContext.Session.Status(), Report: *reportContext.Report}, nil
	}
	if reportContext.Session.Status() != interviewdomain.InterviewStatusInProgress ||
		reportContext.Session.CurrentRound() != interviewdomain.MaxRounds ||
		!hasCompleteInterviewTranscript(reportContext.Messages) {
		return ReportOutput{}, ErrInterviewReportUnavailable
	}

	report, err := service.aiClient.GenerateInterviewReport(ctx, reportContext)
	if err != nil {
		return ReportOutput{}, err
	}
	if err := reportContext.Session.Complete(); err != nil {
		return ReportOutput{}, ErrInterviewReportUnavailable
	}
	if err := service.repository.SaveReport(ctx, reportContext.Session, report); err != nil {
		return ReportOutput{}, err
	}
	return ReportOutput{SessionID: sessionID, Status: reportContext.Session.Status(), Report: report}, nil
}

func hasCompleteInterviewTranscript(messages []interviewdomain.TranscriptMessage) bool {
	interviewerRounds := make(map[int]bool, interviewdomain.MaxRounds)
	candidateRounds := make(map[int]bool, interviewdomain.MaxRounds)
	for _, message := range messages {
		if message.Round < 1 || message.Round > interviewdomain.MaxRounds {
			continue
		}
		switch message.Role {
		case interviewdomain.InterviewMessageRoleInterviewer:
			interviewerRounds[message.Round] = true
		case interviewdomain.InterviewMessageRoleCandidate:
			candidateRounds[message.Round] = message.Feedback != nil
		}
	}
	for round := 1; round <= interviewdomain.MaxRounds; round++ {
		if !interviewerRounds[round] || !candidateRounds[round] {
			return false
		}
	}
	return true
}
