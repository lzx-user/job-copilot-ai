package repository

import (
	"context"

	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

type PlaceholderInterviewRepository struct{}

var _ port.InterviewRepository = (*PlaceholderInterviewRepository)(nil)

func NewPlaceholderInterviewRepository() *PlaceholderInterviewRepository {
	return &PlaceholderInterviewRepository{}
}

func (repository *PlaceholderInterviewRepository) ListOptions(
	context.Context,
	string,
) ([]interviewdomain.InterviewOption, error) {
	return nil, port.ErrRepositoryUnavailable
}

func (repository *PlaceholderInterviewRepository) FindContext(
	context.Context,
	string,
	string,
) (interviewdomain.InterviewContext, error) {
	return interviewdomain.InterviewContext{}, port.ErrRepositoryUnavailable
}

func (repository *PlaceholderInterviewRepository) CreatePendingSession(
	context.Context,
	string,
	string,
) (string, error) {
	return "", port.ErrRepositoryUnavailable
}

func (repository *PlaceholderInterviewRepository) StartSession(
	context.Context,
	interviewdomain.InterviewSession,
	interviewdomain.InterviewMessage,
) error {
	return port.ErrRepositoryUnavailable
}
