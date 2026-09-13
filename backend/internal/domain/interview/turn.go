package interview

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrInvalidInterviewTurn = errors.New("invalid interview turn")

type InterviewFeedback struct {
	score        int
	feedback     string
	strengths    []string
	improvements []string
}

func NewInterviewFeedback(score int, feedback string, strengths, improvements []string) (InterviewFeedback, error) {
	feedback = strings.TrimSpace(feedback)
	if score < 0 || score > 100 || feedback == "" || utf8.RuneCountInString(feedback) > 2000 {
		return InterviewFeedback{}, ErrInvalidInterviewTurn
	}
	strengths, err := validateFeedbackItems(strengths)
	if err != nil {
		return InterviewFeedback{}, err
	}
	improvements, err = validateFeedbackItems(improvements)
	if err != nil {
		return InterviewFeedback{}, err
	}
	return InterviewFeedback{score: score, feedback: feedback, strengths: strengths, improvements: improvements}, nil
}

func validateFeedbackItems(items []string) ([]string, error) {
	if len(items) == 0 || len(items) > 10 {
		return nil, ErrInvalidInterviewTurn
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" || utf8.RuneCountInString(item) > 500 {
			return nil, ErrInvalidInterviewTurn
		}
		result = append(result, item)
	}
	return result, nil
}

type InterviewTurnResult struct {
	feedback     InterviewFeedback
	nextQuestion string
}

func NewInterviewTurnResult(feedback InterviewFeedback, nextQuestion string) (InterviewTurnResult, error) {
	nextQuestion = strings.TrimSpace(nextQuestion)
	if nextQuestion == "" || utf8.RuneCountInString(nextQuestion) > 2000 {
		return InterviewTurnResult{}, ErrInvalidInterviewTurn
	}
	return InterviewTurnResult{feedback: feedback, nextQuestion: nextQuestion}, nil
}

type TranscriptMessage struct {
	ID        string
	Role      InterviewMessageRole
	Round     int
	Content   string
	CreatedAt time.Time
	Feedback  *InterviewFeedback
}

type SessionDetail struct {
	Session     InterviewSession
	CompanyName string
	JobTitle    string
	Messages    []TranscriptMessage
}

type InterviewTurnContext struct {
	Session  InterviewSession
	Context  InterviewContext
	Messages []TranscriptMessage
}

type InterviewTurnPrompt struct {
	Context  InterviewContext
	Messages []TranscriptMessage
	Answer   string
}

func (feedback InterviewFeedback) Score() int       { return feedback.score }
func (feedback InterviewFeedback) Feedback() string { return feedback.feedback }
func (feedback InterviewFeedback) Strengths() []string {
	return append([]string(nil), feedback.strengths...)
}
func (feedback InterviewFeedback) Improvements() []string {
	return append([]string(nil), feedback.improvements...)
}
func (result InterviewTurnResult) Feedback() InterviewFeedback { return result.feedback }
func (result InterviewTurnResult) NextQuestion() string        { return result.nextQuestion }
