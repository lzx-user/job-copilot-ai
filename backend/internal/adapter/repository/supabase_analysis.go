package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	analysisdomain "job-copilot-backend/internal/domain/analysis"
	"job-copilot-backend/internal/port"
)

type SupabaseAnalysisRepository struct {
	endpoint   string
	anonKey    string
	httpClient *http.Client
}

var _ port.AnalysisRepository = (*SupabaseAnalysisRepository)(nil)

func NewSupabaseAnalysisRepository(
	baseURL string,
	anonKey string,
	httpClient *http.Client,
) (*SupabaseAnalysisRepository, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	anonKey = strings.TrimSpace(anonKey)
	if baseURL == "" || anonKey == "" {
		return nil, port.ErrRepositoryUnavailable
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &SupabaseAnalysisRepository{
		endpoint:   baseURL + "/rest/v1/jd_analyses",
		anonKey:    anonKey,
		httpClient: httpClient,
	}, nil
}

type analysisRow struct {
	UserID            string   `json:"user_id"`
	CompanyName       string   `json:"company_name"`
	JobTitle          string   `json:"job_title"`
	JDContent         string   `json:"jd_content"`
	ResumeSummary     string   `json:"resume_summary"`
	Skills            []string `json:"skills"`
	MatchScore        int      `json:"match_score"`
	JobSummary        string   `json:"job_summary"`
	CoreRequirements  []string `json:"core_requirements"`
	MatchedSkills     []string `json:"matched_skills"`
	MissingSkills     []string `json:"missing_skills"`
	ResumeSuggestions []string `json:"resume_suggestions"`
	PreparationTopics []string `json:"preparation_topics"`
	GreetingMessage   string   `json:"greeting_message"`
}

func (repository *SupabaseAnalysisRepository) Save(
	ctx context.Context,
	userID string,
	request analysisdomain.AnalysisRequest,
	result analysisdomain.AnalysisResult,
) (string, error) {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return "", port.ErrUnauthenticated
	}

	payload, err := json.Marshal(analysisRow{
		UserID:            userID,
		CompanyName:       request.CompanyName(),
		JobTitle:          request.JobTitle(),
		JDContent:         request.Description().Content(),
		ResumeSummary:     request.ResumeSummary(),
		Skills:            request.Skills(),
		MatchScore:        result.MatchScore(),
		JobSummary:        result.JobSummary(),
		CoreRequirements:  result.CoreRequirements(),
		MatchedSkills:     result.MatchedSkills(),
		MissingSkills:     result.MissingSkills(),
		ResumeSuggestions: result.ResumeSuggestions(),
		PreparationTopics: result.PreparationTopics(),
		GreetingMessage:   result.GreetingMessage(),
	})
	if err != nil {
		return "", errors.Join(port.ErrRepositoryOperation, err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		repository.endpoint+"?select=id",
		bytes.NewReader(payload),
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
		return "", decodeRepositoryError(httpResponse)
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

func (repository *SupabaseAnalysisRepository) FindByID(
	ctx context.Context,
	userID string,
	analysisID string,
) (analysisdomain.AnalysisResult, error) {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return analysisdomain.AnalysisResult{}, port.ErrUnauthenticated
	}

	query := url.Values{}
	query.Set("id", "eq."+analysisID)
	query.Set("user_id", "eq."+userID)
	query.Set("select", "match_score,job_summary,core_requirements,matched_skills,missing_skills,resume_suggestions,preparation_topics,greeting_message")
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, repository.endpoint+"?"+query.Encode(), nil)
	if err != nil {
		return analysisdomain.AnalysisResult{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	repository.setHeaders(httpRequest, accessToken)

	httpResponse, err := repository.httpClient.Do(httpRequest)
	if err != nil {
		return analysisdomain.AnalysisResult{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, httpResponse.Body)
		return analysisdomain.AnalysisResult{}, fmt.Errorf("%w: status %d", port.ErrRepositoryOperation, httpResponse.StatusCode)
	}

	var rows []struct {
		MatchScore        int      `json:"match_score"`
		JobSummary        string   `json:"job_summary"`
		CoreRequirements  []string `json:"core_requirements"`
		MatchedSkills     []string `json:"matched_skills"`
		MissingSkills     []string `json:"missing_skills"`
		ResumeSuggestions []string `json:"resume_suggestions"`
		PreparationTopics []string `json:"preparation_topics"`
		GreetingMessage   string   `json:"greeting_message"`
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&rows); err != nil || len(rows) != 1 {
		return analysisdomain.AnalysisResult{}, port.ErrRepositoryOperation
	}
	return analysisdomain.NewAnalysisResult(analysisdomain.AnalysisResultParams{
		MatchScore:        rows[0].MatchScore,
		JobSummary:        rows[0].JobSummary,
		CoreRequirements:  rows[0].CoreRequirements,
		MatchedSkills:     rows[0].MatchedSkills,
		MissingSkills:     rows[0].MissingSkills,
		ResumeSuggestions: rows[0].ResumeSuggestions,
		PreparationTopics: rows[0].PreparationTopics,
		GreetingMessage:   rows[0].GreetingMessage,
	})
}

func (repository *SupabaseAnalysisRepository) setHeaders(request *http.Request, accessToken string) {
	request.Header.Set("apikey", repository.anonKey)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
}

func decodeRepositoryError(response *http.Response) error {
	var body struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Hint    string `json:"hint"`
	}
	decoder := json.NewDecoder(io.LimitReader(response.Body, 16<<10))
	if err := decoder.Decode(&body); err != nil {
		return fmt.Errorf("%w: status %d", port.ErrRepositoryOperation, response.StatusCode)
	}
	return fmt.Errorf(
		"%w: status %d code %s message %s hint %s",
		port.ErrRepositoryOperation,
		response.StatusCode,
		body.Code,
		body.Message,
		body.Hint,
	)
}
