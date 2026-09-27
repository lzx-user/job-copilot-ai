package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	dashboarddomain "job-copilot-backend/internal/domain/dashboard"
	"job-copilot-backend/internal/port"
)

type SupabaseDashboardRepository struct {
	endpoint   string
	anonKey    string
	httpClient *http.Client
}

var _ port.DashboardRepository = (*SupabaseDashboardRepository)(nil)

func NewSupabaseDashboardRepository(baseURL, anonKey string, httpClient *http.Client) (*SupabaseDashboardRepository, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	anonKey = strings.TrimSpace(anonKey)
	if baseURL == "" || anonKey == "" {
		return nil, port.ErrRepositoryUnavailable
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &SupabaseDashboardRepository{
		endpoint: baseURL + "/rest/v1/rpc/get_dashboard_summary", anonKey: anonKey, httpClient: httpClient,
	}, nil
}

func (repository *SupabaseDashboardRepository) GetSummary(
	ctx context.Context,
	userID string,
) (dashboarddomain.Summary, error) {
	if strings.TrimSpace(userID) == "" {
		return dashboarddomain.Summary{}, port.ErrRepositoryOperation
	}
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return dashboarddomain.Summary{}, port.ErrUnauthenticated
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, repository.endpoint, bytes.NewReader([]byte(`{}`)))
	if err != nil {
		return dashboarddomain.Summary{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	request.Header.Set("apikey", repository.anonKey)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
	httpResponse, err := repository.httpClient.Do(request)
	if err != nil {
		return dashboarddomain.Summary{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return dashboarddomain.Summary{}, repositoryHTTPError(httpResponse.StatusCode)
	}
	var payload struct {
		JDCount             int `json:"jdCount"`
		InterviewCount      int `json:"interviewCount"`
		WeeklyRecordCount   int `json:"weeklyRecordCount"`
		ProfileCompleteness int `json:"profileCompleteness"`
		RecentRecords       []struct {
			ID          string    `json:"id"`
			Kind        string    `json:"kind"`
			CompanyName string    `json:"companyName"`
			JobTitle    string    `json:"jobTitle"`
			Description string    `json:"description"`
			CreatedAt   time.Time `json:"createdAt"`
		} `json:"recentRecords"`
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&payload); err != nil {
		return dashboarddomain.Summary{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	recentRecords := make([]dashboarddomain.RecentRecord, 0, len(payload.RecentRecords))
	for _, record := range payload.RecentRecords {
		recentRecords = append(recentRecords, dashboarddomain.RecentRecord{
			ID: record.ID, Kind: record.Kind, CompanyName: record.CompanyName,
			JobTitle: record.JobTitle, Description: record.Description, CreatedAt: record.CreatedAt,
		})
	}
	return dashboarddomain.Validate(dashboarddomain.Summary{
		JDCount: payload.JDCount, InterviewCount: payload.InterviewCount,
		WeeklyRecordCount: payload.WeeklyRecordCount, ProfileCompleteness: payload.ProfileCompleteness,
		RecentRecords: recentRecords,
	})
}
