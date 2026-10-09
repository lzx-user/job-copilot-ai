package career

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"unicode/utf8"

	careerdomain "job-copilot-backend/internal/domain/career"
	"job-copilot-backend/internal/port"
)

var (
	ErrMissingDependency = errors.New("missing career dependency")
	ErrInvalidInput      = errors.New("invalid career input")
)

type Service struct{ repository port.CareerRepository }

func NewService(repository port.CareerRepository) (*Service, error) {
	if repository == nil {
		return nil, ErrMissingDependency
	}
	return &Service{repository: repository}, nil
}

func (service *Service) GetWorkspace(ctx context.Context, userID string) (careerdomain.Workspace, error) {
	if strings.TrimSpace(userID) == "" {
		return careerdomain.Workspace{}, ErrInvalidInput
	}
	return service.repository.GetWorkspace(ctx, userID)
}

func (service *Service) CreateResume(ctx context.Context, userID string, input careerdomain.CreateResumeInput) (string, string, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	input.ProjectSummary = strings.TrimSpace(input.ProjectSummary)
	input.Note = strings.TrimSpace(input.Note)
	input.Skills = normalizeList(input.Skills, 30)
	if strings.TrimSpace(userID) == "" || runeLengthOutside(input.Title, 1, 100) ||
		runeLengthOutside(input.Content, 1, 20000) || utf8.RuneCountInString(input.ProjectSummary) > 5000 ||
		utf8.RuneCountInString(input.Note) > 200 {
		return "", "", ErrInvalidInput
	}
	return service.repository.CreateResume(ctx, userID, input)
}

func (service *Service) CreateResumeVersion(ctx context.Context, userID string, input careerdomain.CreateResumeVersionInput) (string, error) {
	input.ResumeID = strings.TrimSpace(input.ResumeID)
	input.Content = strings.TrimSpace(input.Content)
	input.ProjectSummary = strings.TrimSpace(input.ProjectSummary)
	input.Note = strings.TrimSpace(input.Note)
	input.Skills = normalizeList(input.Skills, 30)
	if strings.TrimSpace(userID) == "" || input.ResumeID == "" || runeLengthOutside(input.Content, 1, 20000) ||
		utf8.RuneCountInString(input.ProjectSummary) > 5000 || utf8.RuneCountInString(input.Note) > 200 {
		return "", ErrInvalidInput
	}
	return service.repository.CreateResumeVersion(ctx, userID, input)
}

func (service *Service) CreateJob(ctx context.Context, userID string, input careerdomain.CreateJobInput) (string, string, error) {
	input.CompanyName = strings.TrimSpace(input.CompanyName)
	input.JobTitle = strings.TrimSpace(input.JobTitle)
	input.Location = strings.TrimSpace(input.Location)
	input.SourceURL = strings.TrimSpace(input.SourceURL)
	input.JDContent = strings.TrimSpace(input.JDContent)
	if strings.TrimSpace(userID) == "" || runeLengthOutside(input.CompanyName, 1, 100) ||
		runeLengthOutside(input.JobTitle, 1, 100) || runeLengthOutside(input.JDContent, 200, 8000) ||
		utf8.RuneCountInString(input.Location) > 100 || utf8.RuneCountInString(input.SourceURL) > 2000 ||
		!validOptionalURL(input.SourceURL) {
		return "", "", ErrInvalidInput
	}
	return service.repository.CreateJob(ctx, userID, input)
}

func (service *Service) CreateJDVersion(ctx context.Context, userID string, input careerdomain.CreateJDVersionInput) (string, error) {
	input.JobID = strings.TrimSpace(input.JobID)
	input.JDContent = strings.TrimSpace(input.JDContent)
	if strings.TrimSpace(userID) == "" || input.JobID == "" || runeLengthOutside(input.JDContent, 200, 8000) {
		return "", ErrInvalidInput
	}
	return service.repository.CreateJDVersion(ctx, userID, input)
}

func (service *Service) CreateApplication(ctx context.Context, userID string, input careerdomain.CreateApplicationInput) (string, error) {
	input.JobID = strings.TrimSpace(input.JobID)
	input.ResumeVersionID = strings.TrimSpace(input.ResumeVersionID)
	input.Status = strings.TrimSpace(input.Status)
	input.Source = strings.TrimSpace(input.Source)
	if strings.TrimSpace(userID) == "" || input.JobID == "" || !oneOf(input.Status,
		"planned", "applied", "screening", "interview", "offer", "rejected", "withdrawn", "accepted") ||
		utf8.RuneCountInString(input.Source) > 100 {
		return "", ErrInvalidInput
	}
	return service.repository.CreateApplication(ctx, userID, input)
}

func (service *Service) AddEvent(ctx context.Context, userID string, input careerdomain.CreateEventInput) (string, error) {
	input.ApplicationID = strings.TrimSpace(input.ApplicationID)
	input.EventType = strings.TrimSpace(input.EventType)
	input.Outcome = strings.TrimSpace(input.Outcome)
	input.Notes = strings.TrimSpace(input.Notes)
	if strings.TrimSpace(userID) == "" || input.ApplicationID == "" || !oneOf(input.EventType,
		"planned", "applied", "screening", "written_test", "interview_scheduled", "interview_completed",
		"offer_received", "rejected", "withdrawn", "offer_accepted", "note") || utf8.RuneCountInString(input.Notes) > 5000 {
		return "", ErrInvalidInput
	}
	return service.repository.AddApplicationEvent(ctx, userID, input)
}

func (service *Service) CreateInterview(ctx context.Context, userID string, input careerdomain.CreateInterviewInput) (string, error) {
	input.ApplicationID = strings.TrimSpace(input.ApplicationID)
	input.RoundName = strings.TrimSpace(input.RoundName)
	input.Format = strings.TrimSpace(input.Format)
	input.LocationOrLink = strings.TrimSpace(input.LocationOrLink)
	input.Result = strings.TrimSpace(input.Result)
	input.Notes = strings.TrimSpace(input.Notes)
	if strings.TrimSpace(userID) == "" || input.ApplicationID == "" || runeLengthOutside(input.RoundName, 1, 100) ||
		!oneOf(input.Result, "scheduled", "completed", "passed", "failed", "cancelled", "pending") ||
		(input.DurationMinutes != nil && (*input.DurationMinutes < 1 || *input.DurationMinutes > 1440)) ||
		utf8.RuneCountInString(input.Notes) > 10000 {
		return "", ErrInvalidInput
	}
	return service.repository.CreateRealInterview(ctx, userID, input)
}

func (service *Service) UpdateInterview(ctx context.Context, userID string, input careerdomain.UpdateInterviewInput) error {
	input.InterviewID = strings.TrimSpace(input.InterviewID)
	input.Result = strings.TrimSpace(input.Result)
	input.Notes = strings.TrimSpace(input.Notes)
	if strings.TrimSpace(userID) == "" || input.InterviewID == "" ||
		!oneOf(input.Result, "scheduled", "completed", "passed", "failed", "cancelled", "pending") ||
		utf8.RuneCountInString(input.Notes) > 10000 {
		return ErrInvalidInput
	}
	return service.repository.UpdateRealInterview(ctx, userID, input)
}

func (service *Service) SaveRetrospective(ctx context.Context, userID string, input careerdomain.SaveRetrospectiveInput) error {
	input.InterviewID = strings.TrimSpace(input.InterviewID)
	input.SelfAssessment = strings.TrimSpace(input.SelfAssessment)
	input.AIAnalysis = strings.TrimSpace(input.AIAnalysis)
	input.Questions = normalizeList(input.Questions, 30)
	input.Strengths = normalizeList(input.Strengths, 20)
	input.Weaknesses = normalizeList(input.Weaknesses, 20)
	input.FollowUpActions = normalizeList(input.FollowUpActions, 20)
	if strings.TrimSpace(userID) == "" || input.InterviewID == "" ||
		utf8.RuneCountInString(input.SelfAssessment) > 10000 || utf8.RuneCountInString(input.AIAnalysis) > 10000 {
		return ErrInvalidInput
	}
	return service.repository.SaveRetrospective(ctx, userID, input)
}

func (service *Service) SaveOffer(ctx context.Context, userID string, input careerdomain.SaveOfferInput) error {
	input.ApplicationID = strings.TrimSpace(input.ApplicationID)
	input.Status = strings.TrimSpace(input.Status)
	input.SalarySummary = strings.TrimSpace(input.SalarySummary)
	input.Notes = strings.TrimSpace(input.Notes)
	if strings.TrimSpace(userID) == "" || input.ApplicationID == "" ||
		!oneOf(input.Status, "pending", "accepted", "declined", "expired") ||
		utf8.RuneCountInString(input.SalarySummary) > 500 || utf8.RuneCountInString(input.Notes) > 5000 {
		return ErrInvalidInput
	}
	return service.repository.SaveOffer(ctx, userID, input)
}

func runeLengthOutside(value string, minLength, maxLength int) bool {
	length := utf8.RuneCountInString(value)
	return length < minLength || length > maxLength
}

func validOptionalURL(value string) bool {
	if value == "" {
		return true
	}
	parsed, err := url.ParseRequestURI(value)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https")
}

func oneOf(value string, allowed ...string) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}
	return false
}

func normalizeList(values []string, limit int) []string {
	result := make([]string, 0, min(len(values), limit))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		key := strings.ToLower(value)
		if value == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, value)
		if len(result) == limit {
			break
		}
	}
	return result
}
