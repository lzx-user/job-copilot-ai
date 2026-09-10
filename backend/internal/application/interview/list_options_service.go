package interview

import (
	"context"
	"strings"

	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

type ListInterviewOptionsService struct {
	repository port.InterviewRepository
}

func NewListInterviewOptionsService(repository port.InterviewRepository) (*ListInterviewOptionsService, error) {
	if repository == nil {
		return nil, ErrMissingDependency
	}
	return &ListInterviewOptionsService{repository: repository}, nil
}

func (service *ListInterviewOptionsService) Execute(
	ctx context.Context,
	userID string,
) ([]interviewdomain.InterviewOption, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrInvalidStartInterviewCommand
	}
	return service.repository.ListOptions(ctx, userID)
}
