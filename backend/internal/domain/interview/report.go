package interview

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var ErrInvalidInterviewReport = errors.New("invalid interview report")

// InterviewReport 是完成五轮面试后才能保存的结构化最终报告。
type InterviewReport struct {
	overallScore      int
	technicalScore    int
	expressionScore   int
	projectDepthScore int
	strengths         []string
	weaknesses        []string
	recommendedTopics []string
	answerTips        []string
	summary           string
}

type InterviewReportParams struct {
	OverallScore      int
	TechnicalScore    int
	ExpressionScore   int
	ProjectDepthScore int
	Strengths         []string
	Weaknesses        []string
	RecommendedTopics []string
	AnswerTips        []string
	Summary           string
}

type InterviewReportContext struct {
	Session  InterviewSession
	Context  InterviewContext
	Messages []TranscriptMessage
	Report   *InterviewReport
}

func NewInterviewReport(params InterviewReportParams) (InterviewReport, error) {
	if !isValidReportScore(params.OverallScore) || !isValidReportScore(params.TechnicalScore) ||
		!isValidReportScore(params.ExpressionScore) || !isValidReportScore(params.ProjectDepthScore) {
		return InterviewReport{}, ErrInvalidInterviewReport
	}

	strengths, err := validateReportItems(params.Strengths)
	if err != nil {
		return InterviewReport{}, err
	}
	weaknesses, err := validateReportItems(params.Weaknesses)
	if err != nil {
		return InterviewReport{}, err
	}
	recommendedTopics, err := validateReportItems(params.RecommendedTopics)
	if err != nil {
		return InterviewReport{}, err
	}
	answerTips, err := validateReportItems(params.AnswerTips)
	if err != nil {
		return InterviewReport{}, err
	}
	summary := strings.TrimSpace(params.Summary)
	if summary == "" || utf8.RuneCountInString(summary) > 3000 {
		return InterviewReport{}, ErrInvalidInterviewReport
	}

	return InterviewReport{
		overallScore: params.OverallScore, technicalScore: params.TechnicalScore,
		expressionScore: params.ExpressionScore, projectDepthScore: params.ProjectDepthScore,
		strengths: strengths, weaknesses: weaknesses, recommendedTopics: recommendedTopics,
		answerTips: answerTips, summary: summary,
	}, nil
}

func isValidReportScore(score int) bool { return score >= 0 && score <= 100 }

func validateReportItems(items []string) ([]string, error) {
	if len(items) == 0 || len(items) > 10 {
		return nil, ErrInvalidInterviewReport
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" || utf8.RuneCountInString(item) > 500 {
			return nil, ErrInvalidInterviewReport
		}
		result = append(result, item)
	}
	return result, nil
}

func (report InterviewReport) OverallScore() int      { return report.overallScore }
func (report InterviewReport) TechnicalScore() int    { return report.technicalScore }
func (report InterviewReport) ExpressionScore() int   { return report.expressionScore }
func (report InterviewReport) ProjectDepthScore() int { return report.projectDepthScore }
func (report InterviewReport) Strengths() []string    { return append([]string(nil), report.strengths...) }
func (report InterviewReport) Weaknesses() []string {
	return append([]string(nil), report.weaknesses...)
}
func (report InterviewReport) RecommendedTopics() []string {
	return append([]string(nil), report.recommendedTopics...)
}
func (report InterviewReport) AnswerTips() []string {
	return append([]string(nil), report.answerTips...)
}
func (report InterviewReport) Summary() string { return report.summary }
