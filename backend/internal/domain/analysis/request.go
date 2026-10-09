package analysis

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const (
	maxCompanyNameLength   = 100
	maxJobTitleLength      = 100
	minJDContentLength     = 200
	maxJDContentLength     = 8000
	maxResumeSummaryLength = 5000
	maxSkillsCount         = 30
	maxSkillLength         = 100
)

var ErrInvalidAnalysisRequest = errors.New("invalid analysis request")

// AnalysisRequest 是进入 AI 前已通过业务校验的岗位与候选人上下文。
type AnalysisRequest struct {
	companyName     string
	jobTitle        string
	description     JobDescription
	resumeSummary   string
	skills          []string
	jobID           string
	jobJDVersionID  string
	resumeVersionID string
}

// WithReferences 只接受可选 UUID，使一次分析能追溯到当时使用的岗位、JD 与简历版本。
func (request AnalysisRequest) WithReferences(jobID, jobJDVersionID, resumeVersionID string) (AnalysisRequest, error) {
	jobID = strings.TrimSpace(jobID)
	jobJDVersionID = strings.TrimSpace(jobJDVersionID)
	resumeVersionID = strings.TrimSpace(resumeVersionID)
	if !validOptionalUUID(jobID) || !validOptionalUUID(jobJDVersionID) || !validOptionalUUID(resumeVersionID) {
		return AnalysisRequest{}, ErrInvalidAnalysisRequest
	}
	request.jobID = jobID
	request.jobJDVersionID = jobJDVersionID
	request.resumeVersionID = resumeVersionID
	return request, nil
}

func validOptionalUUID(value string) bool {
	if value == "" {
		return true
	}
	if len(value) != 36 {
		return false
	}
	for index, character := range value {
		if index == 8 || index == 13 || index == 18 || index == 23 {
			if character != '-' {
				return false
			}
			continue
		}
		if !((character >= '0' && character <= '9') || (character >= 'a' && character <= 'f') ||
			(character >= 'A' && character <= 'F')) {
			return false
		}
	}
	return true
}

func NewAnalysisRequest(
	companyName string,
	jobTitle string,
	jdContent string,
	resumeSummary string,
	skills []string,
) (AnalysisRequest, error) {
	companyName = strings.TrimSpace(companyName)
	jobTitle = strings.TrimSpace(jobTitle)
	resumeSummary = strings.TrimSpace(resumeSummary)

	description, err := NewJobDescription(jdContent)
	if err != nil {
		return AnalysisRequest{}, ErrInvalidAnalysisRequest
	}

	normalizedSkills := normalizeSkills(skills)
	if companyName == "" || utf8.RuneCountInString(companyName) > maxCompanyNameLength ||
		jobTitle == "" || utf8.RuneCountInString(jobTitle) > maxJobTitleLength ||
		utf8.RuneCountInString(description.Content()) < minJDContentLength ||
		utf8.RuneCountInString(description.Content()) > maxJDContentLength ||
		resumeSummary == "" || utf8.RuneCountInString(resumeSummary) > maxResumeSummaryLength ||
		len(normalizedSkills) == 0 || len(normalizedSkills) > maxSkillsCount ||
		hasOversizedSkill(normalizedSkills) {
		return AnalysisRequest{}, ErrInvalidAnalysisRequest
	}

	return AnalysisRequest{
		companyName:   companyName,
		jobTitle:      jobTitle,
		description:   description,
		resumeSummary: resumeSummary,
		skills:        normalizedSkills,
	}, nil
}

func hasOversizedSkill(skills []string) bool {
	for _, skill := range skills {
		if utf8.RuneCountInString(skill) > maxSkillLength {
			return true
		}
	}
	return false
}

func normalizeSkills(skills []string) []string {
	seen := make(map[string]struct{}, len(skills))
	normalized := make([]string, 0, len(skills))
	for _, value := range skills {
		skill := strings.TrimSpace(value)
		key := strings.ToLower(skill)
		if skill == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, skill)
	}
	return normalized
}

func (request AnalysisRequest) CompanyName() string         { return request.companyName }
func (request AnalysisRequest) JobTitle() string            { return request.jobTitle }
func (request AnalysisRequest) Description() JobDescription { return request.description }
func (request AnalysisRequest) ResumeSummary() string       { return request.resumeSummary }
func (request AnalysisRequest) JobID() string               { return request.jobID }
func (request AnalysisRequest) JobJDVersionID() string      { return request.jobJDVersionID }
func (request AnalysisRequest) ResumeVersionID() string     { return request.resumeVersionID }

func (request AnalysisRequest) Skills() []string {
	return append([]string(nil), request.skills...)
}
