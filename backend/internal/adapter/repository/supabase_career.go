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

	careerdomain "job-copilot-backend/internal/domain/career"
	"job-copilot-backend/internal/port"
)

type SupabaseCareerRepository struct {
	baseURL    string
	anonKey    string
	httpClient *http.Client
}

var _ port.CareerRepository = (*SupabaseCareerRepository)(nil)

func NewSupabaseCareerRepository(baseURL, anonKey string, httpClient *http.Client) (*SupabaseCareerRepository, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	anonKey = strings.TrimSpace(anonKey)
	if baseURL == "" || anonKey == "" {
		return nil, port.ErrRepositoryUnavailable
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &SupabaseCareerRepository{baseURL: baseURL + "/rest/v1", anonKey: anonKey, httpClient: httpClient}, nil
}

type resumeRow struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type resumeVersionRow struct {
	ID             string    `json:"id"`
	ResumeID       string    `json:"resume_id"`
	VersionNumber  int       `json:"version_number"`
	Content        string    `json:"content"`
	Skills         []string  `json:"skills"`
	ProjectSummary string    `json:"project_summary"`
	Note           string    `json:"note"`
	CreatedAt      time.Time `json:"created_at"`
}

type jobRow struct {
	ID                string    `json:"id"`
	CompanyName       string    `json:"company_name"`
	JobTitle          string    `json:"job_title"`
	Location          string    `json:"location"`
	SourceURL         string    `json:"source_url"`
	RecruitmentStatus string    `json:"recruitment_status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type jdVersionRow struct {
	ID            string    `json:"id"`
	JobID         string    `json:"job_id"`
	VersionNumber int       `json:"version_number"`
	Content       string    `json:"jd_content"`
	CreatedAt     time.Time `json:"created_at"`
}

type applicationRow struct {
	ID              string     `json:"id"`
	JobID           string     `json:"job_id"`
	ResumeVersionID *string    `json:"resume_version_id"`
	Status          string     `json:"status"`
	Source          string     `json:"source"`
	AppliedAt       *time.Time `json:"applied_at"`
	Deadline        *time.Time `json:"deadline"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type eventRow struct {
	ID            string    `json:"id"`
	ApplicationID string    `json:"application_id"`
	EventType     string    `json:"event_type"`
	OccurredAt    time.Time `json:"occurred_at"`
	Outcome       string    `json:"outcome"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

type realInterviewRow struct {
	ID              string     `json:"id"`
	ApplicationID   string     `json:"application_id"`
	RoundName       string     `json:"round_name"`
	ScheduledAt     *time.Time `json:"scheduled_at"`
	DurationMinutes *int       `json:"duration_minutes"`
	Format          string     `json:"format"`
	LocationOrLink  string     `json:"location_or_link"`
	Result          string     `json:"result"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type retrospectiveRow struct {
	InterviewID     string    `json:"interview_id"`
	Questions       []string  `json:"questions"`
	SelfAssessment  string    `json:"self_assessment"`
	Strengths       []string  `json:"strengths"`
	Weaknesses      []string  `json:"weaknesses"`
	FollowUpActions []string  `json:"follow_up_actions"`
	AIAnalysis      string    `json:"ai_analysis"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type offerRow struct {
	ApplicationID string     `json:"application_id"`
	ReceivedAt    time.Time  `json:"received_at"`
	Status        string     `json:"status"`
	Deadline      *time.Time `json:"deadline"`
	SalarySummary string     `json:"salary_summary"`
	Notes         string     `json:"notes"`
	DecidedAt     *time.Time `json:"decided_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (repository *SupabaseCareerRepository) GetWorkspace(ctx context.Context, userID string) (careerdomain.Workspace, error) {
	accessToken, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return careerdomain.Workspace{}, port.ErrUnauthenticated
	}
	if strings.TrimSpace(userID) == "" {
		return careerdomain.Workspace{}, port.ErrRepositoryOperation
	}

	var resumes []resumeRow
	var resumeVersions []resumeVersionRow
	var jobs []jobRow
	var jdVersions []jdVersionRow
	var applications []applicationRow
	var events []eventRow
	var interviews []realInterviewRow
	var retrospectives []retrospectiveRow
	var offers []offerRow

	queries := []struct {
		table       string
		selectValue string
		order       string
		target      any
	}{
		{"resumes", "id,title,is_default,created_at,updated_at", "updated_at.desc", &resumes},
		{"resume_versions", "id,resume_id,version_number,content,skills,project_summary,note,created_at", "created_at.desc", &resumeVersions},
		{"jobs", "id,company_name,job_title,location,source_url,recruitment_status,created_at,updated_at", "updated_at.desc", &jobs},
		{"job_jd_versions", "id,job_id,version_number,jd_content,created_at", "created_at.desc", &jdVersions},
		{"applications", "id,job_id,resume_version_id,status,source,applied_at,deadline,created_at,updated_at", "updated_at.desc", &applications},
		{"application_events", "id,application_id,event_type,occurred_at,outcome,notes,created_at", "occurred_at.desc", &events},
		{"real_interviews", "id,application_id,round_name,scheduled_at,duration_minutes,format,location_or_link,result,notes,created_at,updated_at", "created_at.desc", &interviews},
		{"interview_retrospectives", "interview_id,questions,self_assessment,strengths,weaknesses,follow_up_actions,ai_analysis,created_at,updated_at", "created_at.desc", &retrospectives},
		{"offers", "application_id,received_at,status,deadline,salary_summary,notes,decided_at,created_at,updated_at", "created_at.desc", &offers},
	}
	for _, query := range queries {
		values := url.Values{}
		values.Set("user_id", "eq."+userID)
		values.Set("select", query.selectValue)
		values.Set("order", query.order)
		if err := repository.get(ctx, accessToken, query.table, values, query.target); err != nil {
			return careerdomain.Workspace{}, err
		}
	}

	versionsByResume := make(map[string][]careerdomain.ResumeVersion)
	for _, row := range resumeVersions {
		versionsByResume[row.ResumeID] = append(versionsByResume[row.ResumeID], careerdomain.ResumeVersion{
			ID: row.ID, ResumeID: row.ResumeID, VersionNumber: row.VersionNumber, Content: row.Content,
			Skills: row.Skills, ProjectSummary: row.ProjectSummary, Note: row.Note, CreatedAt: row.CreatedAt,
		})
	}
	jdByJob := make(map[string][]careerdomain.JDVersion)
	for _, row := range jdVersions {
		jdByJob[row.JobID] = append(jdByJob[row.JobID], careerdomain.JDVersion{
			ID: row.ID, JobID: row.JobID, VersionNumber: row.VersionNumber, Content: row.Content, CreatedAt: row.CreatedAt,
		})
	}

	workspace := careerdomain.Workspace{
		Resumes: make([]careerdomain.Resume, 0, len(resumes)), Jobs: make([]careerdomain.Job, 0, len(jobs)),
		Applications: make([]careerdomain.Application, 0, len(applications)), Events: make([]careerdomain.ApplicationEvent, 0, len(events)),
		Interviews: make([]careerdomain.RealInterview, 0, len(interviews)), Retrospectives: make([]careerdomain.InterviewRetrospective, 0, len(retrospectives)),
		Offers: make([]careerdomain.Offer, 0, len(offers)),
	}
	for _, row := range resumes {
		workspace.Resumes = append(workspace.Resumes, careerdomain.Resume{
			ID: row.ID, Title: row.Title, IsDefault: row.IsDefault, CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt, Versions: versionsByResume[row.ID],
		})
	}
	for _, row := range jobs {
		workspace.Jobs = append(workspace.Jobs, careerdomain.Job{
			ID: row.ID, CompanyName: row.CompanyName, JobTitle: row.JobTitle, Location: row.Location,
			SourceURL: row.SourceURL, RecruitmentStatus: row.RecruitmentStatus, CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt, JDVersions: jdByJob[row.ID],
		})
	}
	for _, row := range applications {
		resumeVersionID := ""
		if row.ResumeVersionID != nil {
			resumeVersionID = *row.ResumeVersionID
		}
		workspace.Applications = append(workspace.Applications, careerdomain.Application{
			ID: row.ID, JobID: row.JobID, ResumeVersionID: resumeVersionID, Status: row.Status,
			Source: row.Source, AppliedAt: row.AppliedAt, Deadline: row.Deadline, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		})
	}
	for _, row := range events {
		workspace.Events = append(workspace.Events, careerdomain.ApplicationEvent{
			ID: row.ID, ApplicationID: row.ApplicationID, EventType: row.EventType, OccurredAt: row.OccurredAt,
			Outcome: row.Outcome, Notes: row.Notes, CreatedAt: row.CreatedAt,
		})
	}
	for _, row := range interviews {
		workspace.Interviews = append(workspace.Interviews, careerdomain.RealInterview{
			ID: row.ID, ApplicationID: row.ApplicationID, RoundName: row.RoundName, ScheduledAt: row.ScheduledAt,
			DurationMinutes: row.DurationMinutes, Format: row.Format, LocationOrLink: row.LocationOrLink,
			Result: row.Result, Notes: row.Notes, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		})
	}
	for _, row := range retrospectives {
		workspace.Retrospectives = append(workspace.Retrospectives, careerdomain.InterviewRetrospective{
			InterviewID: row.InterviewID, Questions: row.Questions, SelfAssessment: row.SelfAssessment,
			Strengths: row.Strengths, Weaknesses: row.Weaknesses, FollowUpActions: row.FollowUpActions,
			AIAnalysis: row.AIAnalysis, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		})
	}
	for _, row := range offers {
		workspace.Offers = append(workspace.Offers, careerdomain.Offer{
			ApplicationID: row.ApplicationID, ReceivedAt: row.ReceivedAt, Status: row.Status,
			Deadline: row.Deadline, SalarySummary: row.SalarySummary, Notes: row.Notes,
			DecidedAt: row.DecidedAt, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		})
	}
	return workspace, nil
}

func (repository *SupabaseCareerRepository) CreateResume(ctx context.Context, _ string, input careerdomain.CreateResumeInput) (string, string, error) {
	var result struct {
		ResumeID  string `json:"resumeId"`
		VersionID string `json:"versionId"`
	}
	err := repository.rpc(ctx, "create_resume_with_version", map[string]any{
		"p_title": input.Title, "p_content": input.Content, "p_skills": input.Skills,
		"p_project_summary": input.ProjectSummary, "p_note": input.Note, "p_is_default": input.IsDefault,
	}, &result)
	return result.ResumeID, result.VersionID, err
}

func (repository *SupabaseCareerRepository) CreateResumeVersion(ctx context.Context, _ string, input careerdomain.CreateResumeVersionInput) (string, error) {
	var id string
	err := repository.rpc(ctx, "create_resume_version", map[string]any{
		"p_resume_id": input.ResumeID, "p_content": input.Content, "p_skills": input.Skills,
		"p_project_summary": input.ProjectSummary, "p_note": input.Note,
	}, &id)
	return id, err
}

func (repository *SupabaseCareerRepository) CreateJob(ctx context.Context, _ string, input careerdomain.CreateJobInput) (string, string, error) {
	var result struct {
		JobID       string `json:"jobId"`
		JDVersionID string `json:"jdVersionId"`
	}
	err := repository.rpc(ctx, "create_job_with_jd", map[string]any{
		"p_company_name": input.CompanyName, "p_job_title": input.JobTitle, "p_location": input.Location,
		"p_source_url": input.SourceURL, "p_jd_content": input.JDContent,
	}, &result)
	return result.JobID, result.JDVersionID, err
}

func (repository *SupabaseCareerRepository) CreateJDVersion(ctx context.Context, _ string, input careerdomain.CreateJDVersionInput) (string, error) {
	var id string
	err := repository.rpc(ctx, "create_job_jd_version", map[string]any{
		"p_job_id": input.JobID, "p_jd_content": input.JDContent,
	}, &id)
	return id, err
}

func (repository *SupabaseCareerRepository) CreateApplication(ctx context.Context, _ string, input careerdomain.CreateApplicationInput) (string, error) {
	var id string
	resumeVersionID := any(nil)
	if input.ResumeVersionID != "" {
		resumeVersionID = input.ResumeVersionID
	}
	err := repository.rpc(ctx, "create_application_with_event", map[string]any{
		"p_job_id": input.JobID, "p_resume_version_id": resumeVersionID, "p_status": input.Status,
		"p_source": input.Source, "p_applied_at": input.AppliedAt, "p_deadline": input.Deadline,
	}, &id)
	return id, err
}

func (repository *SupabaseCareerRepository) AddApplicationEvent(ctx context.Context, _ string, input careerdomain.CreateEventInput) (string, error) {
	var id string
	err := repository.rpc(ctx, "add_application_event", map[string]any{
		"p_application_id": input.ApplicationID, "p_event_type": input.EventType,
		"p_occurred_at": input.OccurredAt, "p_outcome": input.Outcome, "p_notes": input.Notes,
	}, &id)
	return id, err
}

func (repository *SupabaseCareerRepository) CreateRealInterview(ctx context.Context, _ string, input careerdomain.CreateInterviewInput) (string, error) {
	var id string
	err := repository.rpc(ctx, "record_real_interview", map[string]any{
		"p_application_id": input.ApplicationID, "p_round_name": input.RoundName,
		"p_scheduled_at": input.ScheduledAt, "p_duration_minutes": input.DurationMinutes,
		"p_format": input.Format, "p_location_or_link": input.LocationOrLink,
		"p_result": input.Result, "p_notes": input.Notes,
	}, &id)
	return id, err
}

func (repository *SupabaseCareerRepository) UpdateRealInterview(ctx context.Context, _ string, input careerdomain.UpdateInterviewInput) error {
	var saved bool
	return repository.rpc(ctx, "update_real_interview_result", map[string]any{
		"p_interview_id": input.InterviewID, "p_result": input.Result, "p_notes": input.Notes,
	}, &saved)
}

func (repository *SupabaseCareerRepository) SaveRetrospective(ctx context.Context, userID string, input careerdomain.SaveRetrospectiveInput) error {
	return repository.write(ctx, http.MethodPost, "interview_retrospectives", "on_conflict=interview_id", map[string]any{
		"interview_id": input.InterviewID, "user_id": userID, "questions": input.Questions,
		"self_assessment": input.SelfAssessment, "strengths": input.Strengths, "weaknesses": input.Weaknesses,
		"follow_up_actions": input.FollowUpActions, "ai_analysis": input.AIAnalysis, "updated_at": time.Now().UTC(),
	}, "resolution=merge-duplicates,return=minimal", nil)
}

func (repository *SupabaseCareerRepository) SaveOffer(ctx context.Context, _ string, input careerdomain.SaveOfferInput) error {
	receivedAt := time.Now().UTC()
	if input.ReceivedAt != nil {
		receivedAt = input.ReceivedAt.UTC()
	}
	var saved bool
	return repository.rpc(ctx, "save_application_offer", map[string]any{
		"p_application_id": input.ApplicationID, "p_received_at": receivedAt,
		"p_status": input.Status, "p_deadline": input.Deadline, "p_salary_summary": input.SalarySummary,
		"p_notes": input.Notes, "p_decided_at": input.DecidedAt,
	}, &saved)
}

func (repository *SupabaseCareerRepository) get(ctx context.Context, token, table string, query url.Values, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, repository.baseURL+"/"+table+"?"+query.Encode(), nil)
	if err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	repository.setHeaders(request, token)
	response, err := repository.httpClient.Do(request)
	if err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return decodeCareerRepositoryError(response)
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	return nil
}

func decodeCareerRepositoryError(response *http.Response) error {
	var body struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 16<<10)).Decode(&body); err != nil {
		return fmt.Errorf("%w: status %d", port.ErrRepositoryOperation, response.StatusCode)
	}
	message := strings.ToLower(body.Message)
	if body.Code == "PGRST202" || body.Code == "PGRST204" || body.Code == "PGRST205" || body.Code == "42P01" ||
		strings.Contains(message, "schema cache") || strings.Contains(message, "could not find the table") {
		return fmt.Errorf("%w: status %d code %s", port.ErrRepositorySchema, response.StatusCode, body.Code)
	}
	return fmt.Errorf("%w: status %d code %s", port.ErrRepositoryOperation, response.StatusCode, body.Code)
}

func (repository *SupabaseCareerRepository) rpc(ctx context.Context, name string, payload any, target any) error {
	return repository.write(ctx, http.MethodPost, "rpc/"+name, "", payload, "", target)
}

func (repository *SupabaseCareerRepository) write(ctx context.Context, method, path, rawQuery string, payload any, prefer string, target any) error {
	token, ok := port.AuthenticatedAccessToken(ctx)
	if !ok {
		return port.ErrUnauthenticated
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	endpoint := repository.baseURL + "/" + path
	if rawQuery != "" {
		endpoint += "?" + rawQuery
	}
	request, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(body))
	if err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	repository.setHeaders(request, token)
	if prefer != "" {
		request.Header.Set("Prefer", prefer)
	}
	response, err := repository.httpClient.Do(request)
	if err != nil {
		return errors.Join(port.ErrRepositoryOperation, err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return decodeCareerRepositoryError(response)
	}
	if target != nil {
		if err := json.NewDecoder(response.Body).Decode(target); err != nil {
			return errors.Join(port.ErrRepositoryOperation, err)
		}
	}
	return nil
}

func (repository *SupabaseCareerRepository) setHeaders(request *http.Request, token string) {
	request.Header.Set("apikey", repository.anonKey)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
}
