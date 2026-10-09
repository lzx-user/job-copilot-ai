package dashboard

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidSummary = errors.New("invalid dashboard summary")

type RecentRecord struct {
	ID          string
	Kind        string
	CompanyName string
	JobTitle    string
	Description string
	CreatedAt   time.Time
}

type Summary struct {
	JDCount                    int
	InterviewCount             int
	ApplicationCount           int
	ActiveApplicationCount     int
	RealInterviewCount         int
	OfferCount                 int
	ApplicationToInterviewRate int
	InterviewToOfferRate       int
	WeeklyRecordCount          int
	ProfileCompleteness        int
	RecentRecords              []RecentRecord
}

func Validate(summary Summary) (Summary, error) {
	if summary.JDCount < 0 || summary.InterviewCount < 0 || summary.ApplicationCount < 0 ||
		summary.ActiveApplicationCount < 0 || summary.RealInterviewCount < 0 || summary.OfferCount < 0 ||
		summary.ApplicationToInterviewRate < 0 || summary.ApplicationToInterviewRate > 100 ||
		summary.InterviewToOfferRate < 0 || summary.InterviewToOfferRate > 100 || summary.WeeklyRecordCount < 0 ||
		summary.ProfileCompleteness < 0 || summary.ProfileCompleteness > 100 ||
		summary.ProfileCompleteness%25 != 0 || len(summary.RecentRecords) > 5 {
		return Summary{}, ErrInvalidSummary
	}
	for _, record := range summary.RecentRecords {
		if strings.TrimSpace(record.ID) == "" ||
			(record.Kind != "jd" && record.Kind != "interview") ||
			strings.TrimSpace(record.CompanyName) == "" || strings.TrimSpace(record.JobTitle) == "" ||
			strings.TrimSpace(record.Description) == "" || record.CreatedAt.IsZero() {
			return Summary{}, ErrInvalidSummary
		}
	}
	return summary, nil
}
