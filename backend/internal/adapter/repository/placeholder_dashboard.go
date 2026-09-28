package repository

import (
	"context"

	dashboarddomain "job-copilot-backend/internal/domain/dashboard"
	"job-copilot-backend/internal/port"
)

type PlaceholderDashboardRepository struct{}

var _ port.DashboardRepository = (*PlaceholderDashboardRepository)(nil)

func (repository *PlaceholderDashboardRepository) GetSummary(
	context.Context,
	string,
) (dashboarddomain.Summary, error) {
	return dashboarddomain.Summary{}, port.ErrRepositoryUnavailable
}
