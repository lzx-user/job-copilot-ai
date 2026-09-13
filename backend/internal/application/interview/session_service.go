package interview

import (
	"context"
	"errors"
	"strings"

	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

var ErrInvalidInterviewSessionQuery = errors.New("invalid interview session query")

type SessionService struct{ repository port.InterviewRepository }

func NewSessionService(repository port.InterviewRepository) (*SessionService, error) {
	if repository == nil {
		return nil, ErrMissingDependency
	}
	return &SessionService{repository: repository}, nil
}

func (service *SessionService) Get(ctx context.Context, userID, sessionID string) (interviewdomain.SessionDetail, error) {
	if strings.TrimSpace(userID) == "" || !isUUID(strings.TrimSpace(sessionID)) {
		return interviewdomain.SessionDetail{}, ErrInvalidInterviewSessionQuery
	}
	return service.repository.FindSessionDetail(ctx, userID, sessionID)
}
