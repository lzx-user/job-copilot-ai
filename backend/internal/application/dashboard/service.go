package dashboard

import (
	"context"
	"errors"
	"strings"

	dashboarddomain "job-copilot-backend/internal/domain/dashboard"
	"job-copilot-backend/internal/port"
)

var (
	ErrMissingDependency = errors.New("missing dashboard dependency")
	ErrInvalidQuery      = errors.New("invalid dashboard query")
)

type Service struct{ repository port.DashboardRepository }

func NewService(repository port.DashboardRepository) (*Service, error) {
	if repository == nil {
		return nil, ErrMissingDependency
	}
	return &Service{repository: repository}, nil
}

func (service *Service) Get(ctx context.Context, userID string) (dashboarddomain.Summary, error) {
	if strings.TrimSpace(userID) == "" {
		return dashboarddomain.Summary{}, ErrInvalidQuery
	}
	return service.repository.GetSummary(ctx, userID)
}
