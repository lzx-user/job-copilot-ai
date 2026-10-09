package repository

import (
	"context"

	careerdomain "job-copilot-backend/internal/domain/career"
	"job-copilot-backend/internal/port"
)

// PlaceholderCareerRepository 让缺少 Supabase 配置时服务仍可启动，并返回可区分的配置错误。
type PlaceholderCareerRepository struct{}

func (repository *PlaceholderCareerRepository) GetWorkspace(context.Context, string) (careerdomain.Workspace, error) {
	return careerdomain.Workspace{}, port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) CreateResume(context.Context, string, careerdomain.CreateResumeInput) (string, string, error) {
	return "", "", port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) CreateResumeVersion(context.Context, string, careerdomain.CreateResumeVersionInput) (string, error) {
	return "", port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) CreateJob(context.Context, string, careerdomain.CreateJobInput) (string, string, error) {
	return "", "", port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) CreateJDVersion(context.Context, string, careerdomain.CreateJDVersionInput) (string, error) {
	return "", port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) CreateApplication(context.Context, string, careerdomain.CreateApplicationInput) (string, error) {
	return "", port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) AddApplicationEvent(context.Context, string, careerdomain.CreateEventInput) (string, error) {
	return "", port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) CreateRealInterview(context.Context, string, careerdomain.CreateInterviewInput) (string, error) {
	return "", port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) UpdateRealInterview(context.Context, string, careerdomain.UpdateInterviewInput) error {
	return port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) SaveRetrospective(context.Context, string, careerdomain.SaveRetrospectiveInput) error {
	return port.ErrRepositoryUnavailable
}

func (repository *PlaceholderCareerRepository) SaveOffer(context.Context, string, careerdomain.SaveOfferInput) error {
	return port.ErrRepositoryUnavailable
}
