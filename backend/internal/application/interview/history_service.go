package interview

import (
	"context"
	"strings"

	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

type HistoryService struct{ repository port.InterviewRepository }

func NewHistoryService(repository port.InterviewRepository) (*HistoryService, error) {
	if repository == nil {
		return nil, ErrMissingDependency
	}
	return &HistoryService{repository: repository}, nil
}

func (service *HistoryService) List(
	ctx context.Context,
	userID string,
) ([]interviewdomain.InterviewHistoryRecord, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidInterviewSessionQuery
	}
	return service.repository.ListSessions(ctx, userID, 50)
}

func (service *HistoryService) Delete(ctx context.Context, userID, sessionID string) error {
	if strings.TrimSpace(userID) == "" || !isUUID(strings.TrimSpace(sessionID)) {
		return ErrInvalidInterviewSessionQuery
	}
	return service.repository.DeleteSession(ctx, userID, sessionID)
}
