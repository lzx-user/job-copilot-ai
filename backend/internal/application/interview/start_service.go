package interview

import (
	"context"
	"errors"
	"strings"

	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

var (
	ErrMissingDependency            = errors.New("interview service dependency is missing")
	ErrInvalidStartInterviewCommand = errors.New("authenticated user ID and analysis ID are required")
)

type StartInterviewService struct {
	aiClient   port.AIClient
	repository port.InterviewRepository
}

type StartInterviewCommand struct {
	AnalysisID string
}

type StartInterviewOutput struct {
	SessionID    string
	Status       interviewdomain.InterviewStatus
	CurrentRound int
	MaxRounds    int
	Question     string
}

func NewStartInterviewService(
	aiClient port.AIClient,
	repository port.InterviewRepository,
) (*StartInterviewService, error) {
	if aiClient == nil || repository == nil {
		return nil, ErrMissingDependency
	}
	return &StartInterviewService{aiClient: aiClient, repository: repository}, nil
}

func (service *StartInterviewService) Execute(
	ctx context.Context,
	userID string,
	command StartInterviewCommand,
) (StartInterviewOutput, error) {
	userID = strings.TrimSpace(userID)
	analysisID := strings.TrimSpace(command.AnalysisID)
	if userID == "" || !isUUID(analysisID) {
		return StartInterviewOutput{}, ErrInvalidStartInterviewCommand
	}

	interviewContext, err := service.repository.FindContext(ctx, userID, analysisID)
	if err != nil {
		return StartInterviewOutput{}, err
	}

	pendingSession, err := interviewdomain.NewInterviewSession(userID, analysisID)
	if err != nil {
		return StartInterviewOutput{}, err
	}
	// 先生成并校验第一题，避免上游失败时留下无法使用的 pending 会话。
	questionContent, err := service.aiClient.GenerateFirstInterviewQuestion(ctx, interviewContext)
	if err != nil {
		return StartInterviewOutput{}, err
	}
	sessionID, err := service.repository.CreatePendingSession(ctx, userID, analysisID)
	if err != nil {
		return StartInterviewOutput{}, err
	}
	pendingSession, err = interviewdomain.RestoreInterviewSession(interviewdomain.InterviewSessionParams{
		ID:           sessionID,
		UserID:       pendingSession.UserID(),
		AnalysisID:   pendingSession.AnalysisID(),
		Status:       pendingSession.Status(),
		CurrentRound: pendingSession.CurrentRound(),
		MaxRounds:    pendingSession.MaxRounds(),
	})
	if err != nil {
		return StartInterviewOutput{}, err
	}

	if err := pendingSession.Start(); err != nil {
		return StartInterviewOutput{}, err
	}
	question, err := interviewdomain.NewInterviewMessage(interviewdomain.InterviewMessageParams{
		SessionID: pendingSession.ID(),
		UserID:    pendingSession.UserID(),
		Role:      interviewdomain.InterviewMessageRoleInterviewer,
		Round:     pendingSession.CurrentRound(),
		Content:   questionContent,
	})
	if err != nil {
		return StartInterviewOutput{}, errors.Join(port.ErrAIInvalidResponse, err)
	}
	if err := service.repository.StartSession(ctx, pendingSession, question); err != nil {
		return StartInterviewOutput{}, err
	}

	return StartInterviewOutput{
		SessionID:    pendingSession.ID(),
		Status:       pendingSession.Status(),
		CurrentRound: pendingSession.CurrentRound(),
		MaxRounds:    pendingSession.MaxRounds(),
		Question:     question.Content(),
	}, nil
}

func isUUID(value string) bool {
	if len(value) != 36 {
		return false
	}
	for index, character := range value {
		if index == 8 || index == 13 || index == 18 || index == 23 {
			if character != '-' {
				return false
			}
			continue
		}
		if !((character >= '0' && character <= '9') ||
			(character >= 'a' && character <= 'f') ||
			(character >= 'A' && character <= 'F')) {
			return false
		}
	}
	return true
}
