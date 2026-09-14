package interview

import "fmt"

type InterviewStatus string

const (
	InterviewStatusPending    InterviewStatus = "pending"
	InterviewStatusInProgress InterviewStatus = "in_progress"
	InterviewStatusCompleted  InterviewStatus = "completed"
)

func NewInterviewStatus(value string) (InterviewStatus, error) {
	status := InterviewStatus(value)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid interview status: %q", value)
	}
	return status, nil
}

func (status InterviewStatus) IsValid() bool {
	switch status {
	case InterviewStatusPending, InterviewStatusInProgress, InterviewStatusCompleted:
		return true
	default:
		return false
	}
}
