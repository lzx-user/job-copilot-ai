package career

import (
	"context"
	"strings"
	"testing"

	careerdomain "job-copilot-backend/internal/domain/career"
)

type careerRepositoryStub struct {
	createdJob careerdomain.CreateJobInput
}

func (stub *careerRepositoryStub) GetWorkspace(context.Context, string) (careerdomain.Workspace, error) {
	return careerdomain.Workspace{}, nil
}
func (stub *careerRepositoryStub) CreateResume(context.Context, string, careerdomain.CreateResumeInput) (string, string, error) {
	return "resume", "version", nil
}
func (stub *careerRepositoryStub) CreateResumeVersion(context.Context, string, careerdomain.CreateResumeVersionInput) (string, error) {
	return "version", nil
}
func (stub *careerRepositoryStub) CreateJob(_ context.Context, _ string, input careerdomain.CreateJobInput) (string, string, error) {
	stub.createdJob = input
	return "job", "jd", nil
}
func (stub *careerRepositoryStub) CreateJDVersion(context.Context, string, careerdomain.CreateJDVersionInput) (string, error) {
	return "jd", nil
}
func (stub *careerRepositoryStub) CreateApplication(context.Context, string, careerdomain.CreateApplicationInput) (string, error) {
	return "application", nil
}
func (stub *careerRepositoryStub) AddApplicationEvent(context.Context, string, careerdomain.CreateEventInput) (string, error) {
	return "event", nil
}
func (stub *careerRepositoryStub) CreateRealInterview(context.Context, string, careerdomain.CreateInterviewInput) (string, error) {
	return "interview", nil
}
func (stub *careerRepositoryStub) UpdateRealInterview(context.Context, string, careerdomain.UpdateInterviewInput) error {
	return nil
}
func (stub *careerRepositoryStub) SaveRetrospective(context.Context, string, careerdomain.SaveRetrospectiveInput) error {
	return nil
}
func (stub *careerRepositoryStub) SaveOffer(context.Context, string, careerdomain.SaveOfferInput) error {
	return nil
}

func TestCreateJobValidatesAndNormalizesInput(t *testing.T) {
	repository := &careerRepositoryStub{}
	service, err := NewService(repository)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	_, _, err = service.CreateJob(context.Background(), "user-id", careerdomain.CreateJobInput{
		CompanyName: "  示例公司 ", JobTitle: " 后端工程师 ",
		JDContent: strings.Repeat("岗位职责与任职要求。", 25),
	})
	if err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}
	if repository.createdJob.CompanyName != "示例公司" || repository.createdJob.JobTitle != "后端工程师" {
		t.Fatalf("input was not normalized: %#v", repository.createdJob)
	}
}

func TestCreateApplicationRejectsUnknownStatus(t *testing.T) {
	service, err := NewService(&careerRepositoryStub{})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	_, err = service.CreateApplication(context.Background(), "user-id", careerdomain.CreateApplicationInput{
		JobID: "job-id", Status: "unknown-status",
	})
	if err != ErrInvalidInput {
		t.Fatalf("CreateApplication() error = %v, want %v", err, ErrInvalidInput)
	}
}
