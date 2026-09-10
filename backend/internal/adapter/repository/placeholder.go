package repository

import (
	"context"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	"job-copilot-backend/internal/port"
)

type PlaceholderAnalysisRepository struct{}

var _ port.AnalysisRepository = (*PlaceholderAnalysisRepository)(nil)

func NewPlaceholderAnalysisRepository() *PlaceholderAnalysisRepository {
	return &PlaceholderAnalysisRepository{}
}

func (repository *PlaceholderAnalysisRepository) Save(
	context.Context,
	string,
	analysisdomain.AnalysisRequest,
	analysisdomain.AnalysisResult,
) (string, error) {
	return "", port.ErrRepositoryUnavailable
}

func (repository *PlaceholderAnalysisRepository) FindByID(
	context.Context,
	string,
	string,
) (analysisdomain.AnalysisResult, error) {
	return analysisdomain.AnalysisResult{}, port.ErrRepositoryUnavailable
}
