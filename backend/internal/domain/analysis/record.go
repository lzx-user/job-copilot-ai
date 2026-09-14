package analysis

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidAnalysisRecord = errors.New("invalid analysis record")

// AnalysisRecord 是历史页可恢复的一次完整 JD 分析。
type AnalysisRecord struct {
	id          string
	companyName string
	jobTitle    string
	createdAt   time.Time
	result      AnalysisResult
}

func NewAnalysisRecord(id, companyName, jobTitle string, createdAt time.Time, result AnalysisResult) (AnalysisRecord, error) {
	id = strings.TrimSpace(id)
	companyName = strings.TrimSpace(companyName)
	jobTitle = strings.TrimSpace(jobTitle)
	if id == "" || companyName == "" || jobTitle == "" || createdAt.IsZero() {
		return AnalysisRecord{}, ErrInvalidAnalysisRecord
	}
	return AnalysisRecord{id: id, companyName: companyName, jobTitle: jobTitle, createdAt: createdAt, result: result}, nil
}

func (record AnalysisRecord) ID() string             { return record.id }
func (record AnalysisRecord) CompanyName() string    { return record.companyName }
func (record AnalysisRecord) JobTitle() string       { return record.jobTitle }
func (record AnalysisRecord) CreatedAt() time.Time   { return record.createdAt }
func (record AnalysisRecord) Result() AnalysisResult { return record.result }
