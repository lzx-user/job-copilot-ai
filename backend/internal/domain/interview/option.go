package interview

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidInterviewOption = errors.New("invalid interview option")

type InterviewOption struct {
	id          string
	companyName string
	jobTitle    string
	matchScore  int
	createdAt   time.Time
}

func NewInterviewOption(
	id string,
	companyName string,
	jobTitle string,
	matchScore int,
	createdAt time.Time,
) (InterviewOption, error) {
	id = strings.TrimSpace(id)
	companyName = strings.TrimSpace(companyName)
	jobTitle = strings.TrimSpace(jobTitle)
	if id == "" || companyName == "" || jobTitle == "" ||
		matchScore < 0 || matchScore > 100 || createdAt.IsZero() {
		return InterviewOption{}, ErrInvalidInterviewOption
	}
	return InterviewOption{
		id:          id,
		companyName: companyName,
		jobTitle:    jobTitle,
		matchScore:  matchScore,
		createdAt:   createdAt,
	}, nil
}

func (option InterviewOption) ID() string           { return option.id }
func (option InterviewOption) CompanyName() string  { return option.companyName }
func (option InterviewOption) JobTitle() string     { return option.jobTitle }
func (option InterviewOption) MatchScore() int      { return option.matchScore }
func (option InterviewOption) CreatedAt() time.Time { return option.createdAt }
