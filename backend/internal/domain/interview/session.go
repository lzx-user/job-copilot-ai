package interview

import (
	"errors"
	"strings"
)

const MaxRounds = 5

var (
	ErrInvalidInterviewSession = errors.New("invalid interview session")
	ErrInterviewAlreadyStarted = errors.New("interview already started")
	ErrInterviewCannotAdvance  = errors.New("interview cannot advance")
	ErrInterviewCannotComplete = errors.New("interview cannot complete")
)

// InterviewSession 维护五轮面试的状态与轮次约束，不依赖 HTTP 或数据库实现。
type InterviewSession struct {
	id           string
	userID       string
	analysisID   string
	status       InterviewStatus
	currentRound int
	maxRounds    int
}

type InterviewSessionParams struct {
	ID           string
	UserID       string
	AnalysisID   string
	Status       InterviewStatus
	CurrentRound int
	MaxRounds    int
}

func NewInterviewSession(userID string, analysisID string) (InterviewSession, error) {
	return restoreInterviewSession(InterviewSessionParams{
		UserID:       userID,
		AnalysisID:   analysisID,
		Status:       InterviewStatusPending,
		CurrentRound: 0,
		MaxRounds:    MaxRounds,
	})
}

// RestoreInterviewSession 只允许把符合领域约束的持久化数据恢复为实体。
func RestoreInterviewSession(params InterviewSessionParams) (InterviewSession, error) {
	if strings.TrimSpace(params.ID) == "" {
		return InterviewSession{}, ErrInvalidInterviewSession
	}
	return restoreInterviewSession(params)
}

func restoreInterviewSession(params InterviewSessionParams) (InterviewSession, error) {
	params.ID = strings.TrimSpace(params.ID)
	params.UserID = strings.TrimSpace(params.UserID)
	params.AnalysisID = strings.TrimSpace(params.AnalysisID)
	if params.UserID == "" || params.AnalysisID == "" || params.MaxRounds != MaxRounds ||
		!isValidSessionState(params.Status, params.CurrentRound) {
		return InterviewSession{}, ErrInvalidInterviewSession
	}

	return InterviewSession{
		id:           params.ID,
		userID:       params.UserID,
		analysisID:   params.AnalysisID,
		status:       params.Status,
		currentRound: params.CurrentRound,
		maxRounds:    params.MaxRounds,
	}, nil
}

func isValidSessionState(status InterviewStatus, currentRound int) bool {
	switch status {
	case InterviewStatusPending:
		return currentRound == 0
	case InterviewStatusInProgress:
		return currentRound >= 1 && currentRound <= MaxRounds
	case InterviewStatusCompleted:
		return currentRound == MaxRounds
	default:
		return false
	}
}

// Start 在第一道题已成功生成、准备持久化时进入第一轮。
func (session *InterviewSession) Start() error {
	if session.status != InterviewStatusPending {
		return ErrInterviewAlreadyStarted
	}
	session.status = InterviewStatusInProgress
	session.currentRound = 1
	return nil
}

// CanAcceptAnswer 明确第五轮仍可提交回答，但 completed 会话不可继续提交。
func (session InterviewSession) CanAcceptAnswer() bool {
	return session.status == InterviewStatusInProgress &&
		session.currentRound >= 1 && session.currentRound <= session.maxRounds
}

// CanContinue 表示当前回答后是否还能进入下一轮。
func (session InterviewSession) CanContinue() bool {
	return session.status == InterviewStatusInProgress && session.currentRound < session.maxRounds
}

func (session *InterviewSession) AdvanceRound() error {
	if !session.CanContinue() {
		return ErrInterviewCannotAdvance
	}
	session.currentRound++
	return nil
}

func (session *InterviewSession) Complete() error {
	if session.status != InterviewStatusInProgress || session.currentRound != session.maxRounds {
		return ErrInterviewCannotComplete
	}
	session.status = InterviewStatusCompleted
	return nil
}

func (session InterviewSession) ID() string              { return session.id }
func (session InterviewSession) UserID() string          { return session.userID }
func (session InterviewSession) AnalysisID() string      { return session.analysisID }
func (session InterviewSession) Status() InterviewStatus { return session.status }
func (session InterviewSession) CurrentRound() int       { return session.currentRound }
func (session InterviewSession) MaxRounds() int          { return session.maxRounds }
