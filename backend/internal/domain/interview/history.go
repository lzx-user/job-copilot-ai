package interview

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidInterviewHistoryRecord = errors.New("invalid interview history record")

// InterviewHistoryRecord 是历史页需要的最小会话摘要。
type InterviewHistoryRecord struct {
	id           string
	companyName  string
	jobTitle     string
	status       InterviewStatus
	currentRound int
	maxRounds    int
	overallScore *int
	createdAt    time.Time
	updatedAt    time.Time
	completedAt  *time.Time
}

type InterviewHistoryRecordParams struct {
	ID           string
	CompanyName  string
	JobTitle     string
	Status       InterviewStatus
	CurrentRound int
	MaxRounds    int
	OverallScore *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CompletedAt  *time.Time
}

func NewInterviewHistoryRecord(params InterviewHistoryRecordParams) (InterviewHistoryRecord, error) {
	params.ID = strings.TrimSpace(params.ID)
	params.CompanyName = strings.TrimSpace(params.CompanyName)
	params.JobTitle = strings.TrimSpace(params.JobTitle)
	if params.ID == "" || params.CompanyName == "" || params.JobTitle == "" ||
		params.MaxRounds != MaxRounds || !isValidSessionState(params.Status, params.CurrentRound) ||
		params.CreatedAt.IsZero() || params.UpdatedAt.IsZero() {
		return InterviewHistoryRecord{}, ErrInvalidInterviewHistoryRecord
	}
	if params.OverallScore != nil && (*params.OverallScore < 0 || *params.OverallScore > 100) {
		return InterviewHistoryRecord{}, ErrInvalidInterviewHistoryRecord
	}
	if params.Status == InterviewStatusCompleted && params.CompletedAt == nil {
		return InterviewHistoryRecord{}, ErrInvalidInterviewHistoryRecord
	}

	return InterviewHistoryRecord{
		id: params.ID, companyName: params.CompanyName, jobTitle: params.JobTitle,
		status: params.Status, currentRound: params.CurrentRound, maxRounds: params.MaxRounds,
		overallScore: params.OverallScore, createdAt: params.CreatedAt, updatedAt: params.UpdatedAt,
		completedAt: params.CompletedAt,
	}, nil
}

func (record InterviewHistoryRecord) ID() string              { return record.id }
func (record InterviewHistoryRecord) CompanyName() string     { return record.companyName }
func (record InterviewHistoryRecord) JobTitle() string        { return record.jobTitle }
func (record InterviewHistoryRecord) Status() InterviewStatus { return record.status }
func (record InterviewHistoryRecord) CurrentRound() int       { return record.currentRound }
func (record InterviewHistoryRecord) MaxRounds() int          { return record.maxRounds }
func (record InterviewHistoryRecord) OverallScore() *int      { return record.overallScore }
func (record InterviewHistoryRecord) CreatedAt() time.Time    { return record.createdAt }
func (record InterviewHistoryRecord) UpdatedAt() time.Time    { return record.updatedAt }
func (record InterviewHistoryRecord) CompletedAt() *time.Time { return record.completedAt }
