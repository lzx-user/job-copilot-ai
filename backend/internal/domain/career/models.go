package career

import "time"

// Workspace 是“简历 -> 岗位 -> 投递 -> 面试 -> Offer”的只读聚合结果。
// 各实体保留独立 ID，前端不能用展示名称代替真实关系。
type Workspace struct {
	Resumes        []Resume                 `json:"resumes"`
	Jobs           []Job                    `json:"jobs"`
	Applications   []Application            `json:"applications"`
	Events         []ApplicationEvent       `json:"events"`
	Interviews     []RealInterview          `json:"interviews"`
	Retrospectives []InterviewRetrospective `json:"retrospectives"`
	Offers         []Offer                  `json:"offers"`
}

type Resume struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	IsDefault bool            `json:"isDefault"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Versions  []ResumeVersion `json:"versions"`
}

type ResumeVersion struct {
	ID             string    `json:"id"`
	ResumeID       string    `json:"resumeId"`
	VersionNumber  int       `json:"versionNumber"`
	Content        string    `json:"content"`
	Skills         []string  `json:"skills"`
	ProjectSummary string    `json:"projectSummary"`
	Note           string    `json:"note"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Job struct {
	ID                string      `json:"id"`
	CompanyName       string      `json:"companyName"`
	JobTitle          string      `json:"jobTitle"`
	Location          string      `json:"location"`
	SourceURL         string      `json:"sourceUrl"`
	RecruitmentStatus string      `json:"recruitmentStatus"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
	JDVersions        []JDVersion `json:"jdVersions"`
}

type JDVersion struct {
	ID            string    `json:"id"`
	JobID         string    `json:"jobId"`
	VersionNumber int       `json:"versionNumber"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Application struct {
	ID              string     `json:"id"`
	JobID           string     `json:"jobId"`
	ResumeVersionID string     `json:"resumeVersionId"`
	Status          string     `json:"status"`
	Source          string     `json:"source"`
	AppliedAt       *time.Time `json:"appliedAt"`
	Deadline        *time.Time `json:"deadline"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type ApplicationEvent struct {
	ID            string    `json:"id"`
	ApplicationID string    `json:"applicationId"`
	EventType     string    `json:"eventType"`
	OccurredAt    time.Time `json:"occurredAt"`
	Outcome       string    `json:"outcome"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"createdAt"`
}

type RealInterview struct {
	ID              string     `json:"id"`
	ApplicationID   string     `json:"applicationId"`
	RoundName       string     `json:"roundName"`
	ScheduledAt     *time.Time `json:"scheduledAt"`
	DurationMinutes *int       `json:"durationMinutes"`
	Format          string     `json:"format"`
	LocationOrLink  string     `json:"locationOrLink"`
	Result          string     `json:"result"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type InterviewRetrospective struct {
	InterviewID     string    `json:"interviewId"`
	Questions       []string  `json:"questions"`
	SelfAssessment  string    `json:"selfAssessment"`
	Strengths       []string  `json:"strengths"`
	Weaknesses      []string  `json:"weaknesses"`
	FollowUpActions []string  `json:"followUpActions"`
	AIAnalysis      string    `json:"aiAnalysis"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Offer struct {
	ApplicationID string     `json:"applicationId"`
	ReceivedAt    time.Time  `json:"receivedAt"`
	Status        string     `json:"status"`
	Deadline      *time.Time `json:"deadline"`
	SalarySummary string     `json:"salarySummary"`
	Notes         string     `json:"notes"`
	DecidedAt     *time.Time `json:"decidedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type CreateResumeInput struct {
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	Skills         []string `json:"skills"`
	ProjectSummary string   `json:"projectSummary"`
	Note           string   `json:"note"`
	IsDefault      bool     `json:"isDefault"`
}

type CreateResumeVersionInput struct {
	ResumeID       string   `json:"resumeId"`
	Content        string   `json:"content"`
	Skills         []string `json:"skills"`
	ProjectSummary string   `json:"projectSummary"`
	Note           string   `json:"note"`
}

type CreateJobInput struct {
	CompanyName string `json:"companyName"`
	JobTitle    string `json:"jobTitle"`
	Location    string `json:"location"`
	SourceURL   string `json:"sourceUrl"`
	JDContent   string `json:"jdContent"`
}

type CreateJDVersionInput struct {
	JobID     string `json:"jobId"`
	JDContent string `json:"jdContent"`
}

type CreateApplicationInput struct {
	JobID           string     `json:"jobId"`
	ResumeVersionID string     `json:"resumeVersionId"`
	Status          string     `json:"status"`
	Source          string     `json:"source"`
	AppliedAt       *time.Time `json:"appliedAt"`
	Deadline        *time.Time `json:"deadline"`
}

type CreateEventInput struct {
	ApplicationID string     `json:"applicationId"`
	EventType     string     `json:"eventType"`
	OccurredAt    *time.Time `json:"occurredAt"`
	Outcome       string     `json:"outcome"`
	Notes         string     `json:"notes"`
}

type CreateInterviewInput struct {
	ApplicationID   string     `json:"applicationId"`
	RoundName       string     `json:"roundName"`
	ScheduledAt     *time.Time `json:"scheduledAt"`
	DurationMinutes *int       `json:"durationMinutes"`
	Format          string     `json:"format"`
	LocationOrLink  string     `json:"locationOrLink"`
	Result          string     `json:"result"`
	Notes           string     `json:"notes"`
}

type UpdateInterviewInput struct {
	InterviewID string `json:"interviewId"`
	Result      string `json:"result"`
	Notes       string `json:"notes"`
}

type SaveRetrospectiveInput struct {
	InterviewID     string   `json:"interviewId"`
	Questions       []string `json:"questions"`
	SelfAssessment  string   `json:"selfAssessment"`
	Strengths       []string `json:"strengths"`
	Weaknesses      []string `json:"weaknesses"`
	FollowUpActions []string `json:"followUpActions"`
	AIAnalysis      string   `json:"aiAnalysis"`
}

type SaveOfferInput struct {
	ApplicationID string     `json:"applicationId"`
	ReceivedAt    *time.Time `json:"receivedAt"`
	Status        string     `json:"status"`
	Deadline      *time.Time `json:"deadline"`
	SalarySummary string     `json:"salarySummary"`
	Notes         string     `json:"notes"`
	DecidedAt     *time.Time `json:"decidedAt"`
}
