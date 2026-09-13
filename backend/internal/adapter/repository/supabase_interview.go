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

func (repository *SupabaseInterviewRepository) FindTurnContext(
	ctx context.Context,
	userID string,
	sessionID string,
) (interviewdomain.InterviewTurnContext, error) {
	session, err := repository.findSession(ctx, userID, sessionID)
	if err != nil {
		return interviewdomain.InterviewTurnContext{}, err
	}
	interviewContext, err := repository.FindContext(ctx, userID, session.AnalysisID())
	if err != nil {
		return interviewdomain.InterviewTurnContext{}, err
	}
	messages, err := repository.findMessages(ctx, userID, sessionID, 8)
	if err != nil {
		return interviewdomain.InterviewTurnContext{}, err
	}
	return interviewdomain.InterviewTurnContext{Session: session, Context: interviewContext, Messages: messages}, nil
}

func (repository *SupabaseInterviewRepository) FindSessionDetail(
	ctx context.Context,
	userID string,
	sessionID string,
) (interviewdomain.SessionDetail, error) {
	session, err := repository.findSession(ctx, userID, sessionID)
	if err != nil {
		return interviewdomain.SessionDetail{}, err
	}
	interviewContext, err := repository.FindContext(ctx, userID, session.AnalysisID())
	if err != nil {
		return interviewdomain.SessionDetail{}, err
	}
	messages, err := repository.findMessages(ctx, userID, sessionID, 20)
	if err != nil {
		return interviewdomain.SessionDetail{}, err
	}
	return interviewdomain.SessionDetail{
		Session: session, CompanyName: interviewContext.CompanyName(),
		JobTitle: interviewContext.JobTitle(), Messages: messages,
	}, nil
}

func (repository *SupabaseInterviewRepository) SaveTurn(
	ctx context.Context,
	session interviewdomain.InterviewSession,
	answer interviewdomain.InterviewMessage,
	result interviewdomain.InterviewTurnResult,
) error {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return port.ErrUnauthenticated
	}
	expectedRound := session.CurrentRound() - 1
	if expectedRound < 1 || expectedRound > 4 || answer.Round() != expectedRound ||
		answer.SessionID() != session.ID() || answer.UserID() != session.UserID() ||
		answer.Role() != interviewdomain.InterviewMessageRoleCandidate ||
		session.Status() != interviewdomain.InterviewStatusInProgress {
		return port.ErrRepositoryOperation
	}
	feedback := result.Feedback()
	body, err := json.Marshal(struct {
		SessionID     string   `json:"p_session_id"`
		ExpectedRound int      `json:"p_expected_round"`
		Answer        string   `json:"p_answer"`
		Score         int      `json:"p_score"`
		Feedback      string   `json:"p_feedback"`
		Strengths     []string `json:"p_strengths"`
		Improvements  []string `json:"p_improvements"`
		NextQuestion  string   `json:"p_next_question"`
	}{session.ID(), expectedRound, answer.Content(), feedback.Score(), feedback.Feedback(), feedback.Strengths(), feedback.Improvements(), result.NextQuestion()})
	if err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, repository.baseURL+"/rpc/submit_interview_turn", bytes.NewReader(body))
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
		var errorBody struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(httpResponse.Body).Decode(&errorBody)
		if strings.Contains(errorBody.Message, "interview turn conflict") || httpResponse.StatusCode == http.StatusConflict {
			return port.ErrRepositoryConflict
		}
		return repositoryHTTPError(httpResponse.StatusCode)
	}
	return nil
}

func (repository *SupabaseInterviewRepository) findSession(
	ctx context.Context,
	userID, sessionID string,
) (interviewdomain.InterviewSession, error) {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return interviewdomain.InterviewSession{}, port.ErrUnauthenticated
	}
	query := url.Values{}
	query.Set("id", "eq."+sessionID)
	query.Set("user_id", "eq."+userID)
	query.Set("select", "id,user_id,jd_analysis_id,status,current_round,max_rounds")
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, repository.baseURL+"/interview_sessions?"+query.Encode(), nil)
	if err != nil {
		return interviewdomain.InterviewSession{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	repository.setHeaders(httpRequest, accessToken)
	httpResponse, err := repository.httpClient.Do(httpRequest)
	if err != nil {
		return interviewdomain.InterviewSession{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return interviewdomain.InterviewSession{}, repositoryHTTPError(httpResponse.StatusCode)
	}
	// PostgREST 的 snake_case 字段需要显式标签。
	var rawRows []struct {
		ID           string                          `json:"id"`
		UserID       string                          `json:"user_id"`
		AnalysisID   string                          `json:"jd_analysis_id"`
		Status       interviewdomain.InterviewStatus `json:"status"`
		CurrentRound int                             `json:"current_round"`
		MaxRounds    int                             `json:"max_rounds"`
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&rawRows); err != nil {
		return interviewdomain.InterviewSession{}, errors.Join(port.ErrRepositoryOperation, err)
	}
	if len(rawRows) == 0 {
		return interviewdomain.InterviewSession{}, port.ErrRepositoryNotFound
	}
	if len(rawRows) != 1 {
		return interviewdomain.InterviewSession{}, port.ErrRepositoryOperation
	}
	row := rawRows[0]
	return interviewdomain.RestoreInterviewSession(interviewdomain.InterviewSessionParams{
		ID: row.ID, UserID: row.UserID, AnalysisID: row.AnalysisID, Status: row.Status,
		CurrentRound: row.CurrentRound, MaxRounds: row.MaxRounds,
	})
}

func (repository *SupabaseInterviewRepository) findMessages(
	ctx context.Context,
	userID, sessionID string,
	limit int,
) ([]interviewdomain.TranscriptMessage, error) {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return nil, port.ErrUnauthenticated
	}
	query := url.Values{}
	query.Set("session_id", "eq."+sessionID)
	query.Set("user_id", "eq."+userID)
	query.Set("select", "id,role,round,content,created_at,score,feedback,strengths,improvements")
	query.Set("order", "created_at.desc")
	query.Set("limit", fmt.Sprint(limit))
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, repository.baseURL+"/interview_messages?"+query.Encode(), nil)
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
		ID           string                               `json:"id"`
		Role         interviewdomain.InterviewMessageRole `json:"role"`
		Round        int                                  `json:"round"`
		Content      string                               `json:"content"`
		CreatedAt    time.Time                            `json:"created_at"`
		Score        *int                                 `json:"score"`
		Feedback     *string                              `json:"feedback"`
		Strengths    []string                             `json:"strengths"`
		Improvements []string                             `json:"improvements"`
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&rows); err != nil {
		return nil, errors.Join(port.ErrRepositoryOperation, err)
	}
	messages := make([]interviewdomain.TranscriptMessage, 0, len(rows))
	for index := len(rows) - 1; index >= 0; index-- {
		row := rows[index]
		var feedback *interviewdomain.InterviewFeedback
		if row.Score != nil && row.Feedback != nil {
			value, err := interviewdomain.NewInterviewFeedback(*row.Score, *row.Feedback, row.Strengths, row.Improvements)
			if err != nil {
				return nil, errors.Join(port.ErrRepositoryOperation, err)
			}
			feedback = &value
		}
		messages = append(messages, interviewdomain.TranscriptMessage{
			ID: row.ID, Role: row.Role, Round: row.Round, Content: row.Content,
			CreatedAt: row.CreatedAt, Feedback: feedback,
		})
	}
	return messages, nil
}

func (repository *SupabaseInterviewRepository) setHeaders(request *http.Request, accessToken string) {
	request.Header.Set("apikey", repository.anonKey)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
}

func repositoryHTTPError(status int) error {
	return fmt.Errorf("%w: status %d", port.ErrRepositoryOperation, status)
}
