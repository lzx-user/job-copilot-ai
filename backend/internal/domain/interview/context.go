package interview

import (
	"errors"
	"strings"
)

var ErrInvalidInterviewContext = errors.New("invalid interview context")

// InterviewContext 是生成面试题所需的已持久化 JD 与候选人上下文。
type InterviewContext struct {
	companyName       string
	jobTitle          string
	jdContent         string
	resumeSummary     string
	skills            []string
	coreRequirements  []string
	preparationTopics []string
}

type InterviewContextParams struct {
	CompanyName       string
	JobTitle          string
	JDContent         string
	ResumeSummary     string
	Skills            []string
	CoreRequirements  []string
	PreparationTopics []string
}

func NewInterviewContext(params InterviewContextParams) (InterviewContext, error) {
	params.CompanyName = strings.TrimSpace(params.CompanyName)
	params.JobTitle = strings.TrimSpace(params.JobTitle)
	params.JDContent = strings.TrimSpace(params.JDContent)
	params.ResumeSummary = strings.TrimSpace(params.ResumeSummary)
	params.Skills = normalizeContextItems(params.Skills)
	params.CoreRequirements = normalizeContextItems(params.CoreRequirements)
	params.PreparationTopics = normalizeContextItems(params.PreparationTopics)

	if params.CompanyName == "" || params.JobTitle == "" || params.JDContent == "" ||
		params.ResumeSummary == "" || len(params.Skills) == 0 ||
		len(params.CoreRequirements) == 0 || len(params.PreparationTopics) == 0 {
		return InterviewContext{}, ErrInvalidInterviewContext
	}

	return InterviewContext{
		companyName:       params.CompanyName,
		jobTitle:          params.JobTitle,
		jdContent:         params.JDContent,
		resumeSummary:     params.ResumeSummary,
		skills:            params.Skills,
		coreRequirements:  params.CoreRequirements,
		preparationTopics: params.PreparationTopics,
	}, nil
}

func normalizeContextItems(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			normalized = append(normalized, value)
		}
	}
	return normalized
}

func (context InterviewContext) CompanyName() string   { return context.companyName }
func (context InterviewContext) JobTitle() string      { return context.jobTitle }
func (context InterviewContext) JDContent() string     { return context.jdContent }
func (context InterviewContext) ResumeSummary() string { return context.resumeSummary }
func (context InterviewContext) Skills() []string {
	return append([]string(nil), context.skills...)
}
func (context InterviewContext) CoreRequirements() []string {
	return append([]string(nil), context.coreRequirements...)
}
func (context InterviewContext) PreparationTopics() []string {
	return append([]string(nil), context.preparationTopics...)
}
