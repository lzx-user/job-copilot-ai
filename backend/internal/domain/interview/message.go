package interview

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const maxMessageLength = 10000

type InterviewMessageRole string

const (
	InterviewMessageRoleInterviewer InterviewMessageRole = "interviewer"
	InterviewMessageRoleCandidate   InterviewMessageRole = "candidate"
)

var ErrInvalidInterviewMessage = errors.New("invalid interview message")

// InterviewMessage 表示某一轮中可持久化的一条面试官或候选人消息。
type InterviewMessage struct {
	id        string
	sessionID string
	userID    string
	role      InterviewMessageRole
	round     int
	content   string
}

type InterviewMessageParams struct {
	ID        string
	SessionID string
	UserID    string
	Role      InterviewMessageRole
	Round     int
	Content   string
}

func NewInterviewMessage(params InterviewMessageParams) (InterviewMessage, error) {
	params.ID = strings.TrimSpace(params.ID)
	params.SessionID = strings.TrimSpace(params.SessionID)
	params.UserID = strings.TrimSpace(params.UserID)
	params.Content = strings.TrimSpace(params.Content)

	if params.SessionID == "" || params.UserID == "" ||
		!params.Role.IsValid() || params.Round < 1 || params.Round > MaxRounds ||
		params.Content == "" || utf8.RuneCountInString(params.Content) > maxMessageLength {
		return InterviewMessage{}, ErrInvalidInterviewMessage
	}

	return InterviewMessage{
		id:        params.ID,
		sessionID: params.SessionID,
		userID:    params.UserID,
		role:      params.Role,
		round:     params.Round,
		content:   params.Content,
	}, nil
}

func (role InterviewMessageRole) IsValid() bool {
	return role == InterviewMessageRoleInterviewer || role == InterviewMessageRoleCandidate
}

func (message InterviewMessage) ID() string                 { return message.id }
func (message InterviewMessage) SessionID() string          { return message.sessionID }
func (message InterviewMessage) UserID() string             { return message.userID }
func (message InterviewMessage) Role() InterviewMessageRole { return message.role }
func (message InterviewMessage) Round() int                 { return message.round }
func (message InterviewMessage) Content() string            { return message.content }
