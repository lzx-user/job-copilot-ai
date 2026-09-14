package analysis

import (
	"context"
	"errors"
	"strings"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	"job-copilot-backend/internal/port"
)

var ErrInvalidAnalysisHistoryQuery = errors.New("invalid analysis history query")

type HistoryService struct {
	repository port.AnalysisRepository
}

func NewHistoryService(repository port.AnalysisRepository) (*HistoryService, error) {
	if repository == nil {
		return nil, ErrMissingDependency
	}
	return &HistoryService{repository: repository}, nil
}

func (service *HistoryService) List(ctx context.Context, userID string) ([]analysisdomain.AnalysisRecord, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidAnalysisHistoryQuery
	}
	return service.repository.List(ctx, userID, 50)
}

func (service *HistoryService) Get(ctx context.Context, userID, analysisID string) (analysisdomain.AnalysisRecord, error) {
	if strings.TrimSpace(userID) == "" || !isUUID(strings.TrimSpace(analysisID)) {
		return analysisdomain.AnalysisRecord{}, ErrInvalidAnalysisHistoryQuery
	}
	return service.repository.FindRecordByID(ctx, userID, analysisID)
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
