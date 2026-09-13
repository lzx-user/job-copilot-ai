package interview

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

var (
	ErrInvalidInterviewTurnCommand = errors.New("invalid interview turn command")
	ErrInterviewTurnUnavailable    = errors.New("interview cannot accept another answer in this stage")
)

type TurnService struct {
	aiClient   port.AIClient
	repository port.InterviewRepository
}

type TurnCommand struct {
	SessionID string
	Answer    string
}

type TurnOutput struct {
	SessionID    string
	CurrentRound int
	MaxRounds    int
	Result       interviewdomain.InterviewTurnResult
}

func NewTurnService(aiClient port.AIClient, repository port.InterviewRepository) (*TurnService, error) {
	if aiClient == nil || repository == nil {
		return nil, ErrMissingDependency
	}
	return &TurnService{aiClient: aiClient, repository: repository}, nil
}

func (service *TurnService) Execute(ctx context.Context, userID string, command TurnCommand) (TurnOutput, error) {
	userID = strings.TrimSpace(userID)
	sessionID := strings.TrimSpace(command.SessionID)
	answerContent := strings.TrimSpace(command.Answer)
	if userID == "" || !isUUID(sessionID) || answerContent == "" || utf8.RuneCountInString(answerContent) > 10000 {
		return TurnOutput{}, ErrInvalidInterviewTurnCommand
	}
	turnContext, err := service.repository.FindTurnContext(ctx, userID, sessionID)
	if err != nil {
		return TurnOutput{}, err
	}
	// 8.26 只推进需要生成下一题的第 1～4 轮；第 5 轮最终报告在下一阶段完成。
	if !turnContext.Session.CanContinue() {
		return TurnOutput{}, ErrInterviewTurnUnavailable
	}
	result, err := service.aiClient.EvaluateInterviewAnswer(ctx, interviewdomain.InterviewTurnPrompt{
		Context: turnContext.Context, Messages: turnContext.Messages, Answer: answerContent,
	})
	if err != nil {
		return TurnOutput{}, err
	}
	answeredRound := turnContext.Session.CurrentRound()
	answer, err := interviewdomain.NewInterviewMessage(interviewdomain.InterviewMessageParams{
		SessionID: sessionID, UserID: userID, Role: interviewdomain.InterviewMessageRoleCandidate,
		Round: answeredRound, Content: answerContent,
	})
	if err != nil {
		return TurnOutput{}, err
	}
	if err := turnContext.Session.AdvanceRound(); err != nil {
		return TurnOutput{}, ErrInterviewTurnUnavailable
	}
	if err := service.repository.SaveTurn(ctx, turnContext.Session, answer, result); err != nil {
		return TurnOutput{}, err
	}
	return TurnOutput{
		SessionID: sessionID, CurrentRound: turnContext.Session.CurrentRound(),
		MaxRounds: turnContext.Session.MaxRounds(), Result: result,
	}, nil
}
