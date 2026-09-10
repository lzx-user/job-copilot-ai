package analysis

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const (
	maxResultTextLength = 2000
	maxResultListItems  = 20
	maxResultItemLength = 500
)

var ErrInvalidAnalysisResult = errors.New("invalid analysis result")

// AnalysisResult 只保存已经通过业务校验的 AI 分析结果。
type AnalysisResult struct {
	matchScore        MatchScore
	jobSummary        string
	coreRequirements  []string
	matchedSkills     []string
	missingSkills     []string
	resumeSuggestions []string
	preparationTopics []string
	greetingMessage   string
}

type AnalysisResultParams struct {
	MatchScore        int
	JobSummary        string
	CoreRequirements  []string
	MatchedSkills     []string
	MissingSkills     []string
	ResumeSuggestions []string
	PreparationTopics []string
	GreetingMessage   string
}

func NewAnalysisResult(params AnalysisResultParams) (AnalysisResult, error) {
	score, err := NewMatchScore(params.MatchScore)
	if err != nil {
		return AnalysisResult{}, errors.Join(ErrInvalidAnalysisResult, err)
	}

	jobSummary := strings.TrimSpace(params.JobSummary)
	greetingMessage := strings.TrimSpace(params.GreetingMessage)
	coreRequirements, err := validateResultList(params.CoreRequirements, true)
	if err != nil {
		return AnalysisResult{}, err
	}
	matchedSkills, err := validateResultList(params.MatchedSkills, false)
	if err != nil {
		return AnalysisResult{}, err
	}
	missingSkills, err := validateResultList(params.MissingSkills, false)
	if err != nil {
		return AnalysisResult{}, err
	}
	resumeSuggestions, err := validateResultList(params.ResumeSuggestions, true)
	if err != nil {
		return AnalysisResult{}, err
	}
	preparationTopics, err := validateResultList(params.PreparationTopics, true)
	if err != nil {
		return AnalysisResult{}, err
	}

	if jobSummary == "" || greetingMessage == "" ||
		utf8.RuneCountInString(jobSummary) > maxResultTextLength ||
		utf8.RuneCountInString(greetingMessage) > maxResultTextLength {
		return AnalysisResult{}, ErrInvalidAnalysisResult
	}

	return AnalysisResult{
		matchScore:        score,
		jobSummary:        jobSummary,
		coreRequirements:  coreRequirements,
		matchedSkills:     matchedSkills,
		missingSkills:     missingSkills,
		resumeSuggestions: resumeSuggestions,
		preparationTopics: preparationTopics,
		greetingMessage:   greetingMessage,
	}, nil
}

func validateResultList(values []string, requireItem bool) ([]string, error) {
	if len(values) > maxResultListItems || requireItem && len(values) == 0 {
		return nil, ErrInvalidAnalysisResult
	}

	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || utf8.RuneCountInString(value) > maxResultItemLength {
			return nil, ErrInvalidAnalysisResult
		}
		normalized = append(normalized, value)
	}
	return normalized, nil
}

func (result AnalysisResult) MatchScore() int         { return result.matchScore.Value() }
func (result AnalysisResult) JobSummary() string      { return result.jobSummary }
func (result AnalysisResult) GreetingMessage() string { return result.greetingMessage }
func (result AnalysisResult) CoreRequirements() []string {
	return cloneStrings(result.coreRequirements)
}
func (result AnalysisResult) MatchedSkills() []string { return cloneStrings(result.matchedSkills) }
func (result AnalysisResult) MissingSkills() []string { return cloneStrings(result.missingSkills) }
func (result AnalysisResult) ResumeSuggestions() []string {
	return cloneStrings(result.resumeSuggestions)
}
func (result AnalysisResult) PreparationTopics() []string {
	return cloneStrings(result.preparationTopics)
}

func cloneStrings(values []string) []string {
	return append([]string(nil), values...)
}
