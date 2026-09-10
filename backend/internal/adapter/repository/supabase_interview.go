package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	interviewdomain "job-copilot-backend/internal/domain/interview"
	"job-copilot-backend/internal/port"
)

type SupabaseInterviewRepository struct {
	baseURL    string
	anonKey    string
	httpClient *http.Client
}

var _ port.InterviewRepository = (*SupabaseInterviewRepository)(nil)

func NewSupabaseInterviewRepository(
	baseURL string,
	anonKey string,
	httpClient *http.Client,
) (*SupabaseInterviewRepository, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	anonKey = strings.TrimSpace(anonKey)
	if baseURL == "" || anonKey == "" {
		return nil, port.ErrRepositoryUnavailable
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &SupabaseInterviewRepository{
		baseURL:    baseURL + "/rest/v1",
		anonKey:    anonKey,
		httpClient: httpClient,
	}, nil
}

func (repository *SupabaseInterviewRepository) ListOptions(
	ctx context.Context,
	userID string,
) ([]interviewdomain.InterviewOption, error) {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return nil, port.ErrUnauthenticated
	}

	query := url.Values{}
	query.Set("user_id", "eq."+userID)
	query.Set("select", "id,company_name,job_title,match_score,created_at")
	query.Set("order", "created_at.desc")
	query.Set("limit", "20")
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		repository.baseURL+"/jd_analyses?"+query.Encode(),
		nil,
	)
	if err != nil {
		return nil, errors.Join(port.ErrRepositoryOperation, err)
	}
	repository.setHeaders(httpRequest, accessToken)

	httpResponse, err := repository.httpClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Join(port.ErrRepositoryOperation, err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return nil, repositoryHTTPError(httpResponse.StatusCode)
	}

	var rows []struct {
		ID          string    `json:"id"`
		CompanyName string    `json:"company_name"`
		JobTitle    string    `json:"job_title"`
		MatchScore  int       `json:"match_score"`
		CreatedAt   time.Time `json:"created_at"`
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&rows); err != nil {
		return nil, errors.Join(port.ErrRepositoryOperation, err)
	}

	options := make([]interviewdomain.InterviewOption, 0, len(rows))
	for _, row := range rows {
		option, err := interviewdomain.NewInterviewOption(
			row.ID,
			row.CompanyName,
			row.JobTitle,
			row.MatchScore,
			row.CreatedAt,
		)
		if err != nil {
			return nil, errors.Join(port.ErrRepositoryOperation, err)
		}
		options = append(options, option)
	}
	return options, nil
}

func (repository *SupabaseInterviewRepository) FindContext(
	ctx context.Context,
	userID string,
	analysisID string,
) (interviewdomain.InterviewContext, error) {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return interviewdomain.InterviewContext{}, port.ErrUnauthenticated
	}

	query := url.Values{}
	query.Set("id", "eq."+analysisID)
	query.Set("user_id", "eq."+userID)
	query.Set("select", "company_name,job_title,jd_content,resume_summary,skills,core_requirements,preparation_topics")
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		repository.baseURL+"/jd_analyses?"+query.Encode(),
		nil,
	)
	if err != nil {
		return interviewdomain.InterviewContext{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	repository.setHeaders(httpRequest, accessToken)

	httpResponse, err := repository.httpClient.Do(httpRequest)
	if err != nil {
		return interviewdomain.InterviewContext{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return interviewdomain.InterviewContext{}, repositoryHTTPError(httpResponse.StatusCode)
	}

	var rows []struct {
		CompanyName       string   `json:"company_name"`
		JobTitle          string   `json:"job_title"`
		JDContent         string   `json:"jd_content"`
		ResumeSummary     string   `json:"resume_summary"`
		Skills            []string `json:"skills"`
		CoreRequirements  []string `json:"core_requirements"`
		PreparationTopics []string `json:"preparation_topics"`
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&rows); err != nil {
		return interviewdomain.InterviewContext{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	if len(rows) == 0 {
		return interviewdomain.InterviewContext{}, port.ErrRepositoryNotFound
	}
	if len(rows) != 1 {
		return interviewdomain.InterviewContext{}, port.ErrRepositoryOperation
	}

	return interviewdomain.NewInterviewContext(interviewdomain.InterviewContextParams{
		CompanyName:       rows[0].CompanyName,
		JobTitle:          rows[0].JobTitle,
		JDContent:         rows[0].JDContent,
		ResumeSummary:     rows[0].ResumeSummary,
		Skills:            rows[0].Skills,
		CoreRequirements:  rows[0].CoreRequirements,
		PreparationTopics: rows[0].PreparationTopics,
	})
}

func (repository *SupabaseInterviewRepository) CreatePendingSession(
	ctx context.Context,
	userID string,
	analysisID string,
) (string, error) {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return "", port.ErrUnauthenticated
	}

	body, err := json.Marshal(struct {
		UserID     string `json:"user_id"`
		AnalysisID string `json:"jd_analysis_id"`
	}{UserID: userID, AnalysisID: analysisID})
	if err != nil {
		return "", errors.Join(port.ErrRepositoryOperation, err)
	}
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		repository.baseURL+"/interview_sessions?select=id",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", errors.Join(port.ErrRepositoryOperation, err)
	}
	repository.setHeaders(httpRequest, accessToken)
	httpRequest.Header.Set("Prefer", "return=representation")

	httpResponse, err := repository.httpClient.Do(httpRequest)
	if err != nil {
		return "", errors.Join(port.ErrRepositoryOperation, err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return "", repositoryHTTPError(httpResponse.StatusCode)
	}

	var rows []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&rows); err != nil ||
		len(rows) != 1 || strings.TrimSpace(rows[0].ID) == "" {
		return "", port.ErrRepositoryOperation
	}
	return rows[0].ID, nil
}

func (repository *SupabaseInterviewRepository) StartSession(
	ctx context.Context,
	session interviewdomain.InterviewSession,
	question interviewdomain.InterviewMessage,
) error {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return port.ErrUnauthenticated
	}
	if session.ID() == "" || session.Status() != interviewdomain.InterviewStatusInProgress ||
		session.CurrentRound() != 1 || question.SessionID() != session.ID() ||
		question.UserID() != session.UserID() || question.Role() != interviewdomain.InterviewMessageRoleInterviewer ||
		question.Round() != 1 {
		return port.ErrRepositoryOperation
	}

	body, err := json.Marshal(struct {
		SessionID string `json:"p_session_id"`
		Question  string `json:"p_question"`
	}{SessionID: session.ID(), Question: question.Content()})
	if err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		repository.baseURL+"/rpc/start_interview_session",
		bytes.NewReader(body),
	)
	if err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	repository.setHeaders(httpRequest, accessToken)

	httpResponse, err := repository.httpClient.Do(httpRequest)
	if err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return repositoryHTTPError(httpResponse.StatusCode)
	}
	return nil
}

func (repository *SupabaseInterviewRepository) setHeaders(request *http.Request, accessToken string) {
	request.Header.Set("apikey", repository.anonKey)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
}

func repositoryHTTPError(status int) error {
	return fmt.Errorf("%w: status %d", port.ErrRepositoryOperation, status)
}
