export interface ResumeVersion {
  id: string
  resumeId: string
  versionNumber: number
  content: string
  skills: string[]
  projectSummary: string
  note: string
  createdAt: string
}

export interface Resume {
  id: string
  title: string
  isDefault: boolean
  createdAt: string
  updatedAt: string
  versions: ResumeVersion[]
}

export interface JDVersion {
  id: string
  jobId: string
  versionNumber: number
  content: string
  createdAt: string
}

export interface Job {
  id: string
  companyName: string
  jobTitle: string
  location: string
  sourceUrl: string
  recruitmentStatus: 'open' | 'closed' | 'unknown'
  createdAt: string
  updatedAt: string
  jdVersions: JDVersion[]
}

export type ApplicationStatus = 'planned' | 'applied' | 'screening' | 'interview' | 'offer' | 'rejected' | 'withdrawn' | 'accepted'

export interface Application {
  id: string
  jobId: string
  resumeVersionId: string
  status: ApplicationStatus
  source: string
  appliedAt: string | null
  deadline: string | null
  createdAt: string
  updatedAt: string
}

export interface ApplicationEvent {
  id: string
  applicationId: string
  eventType: string
  occurredAt: string
  outcome: string
  notes: string
  createdAt: string
}

export interface RealInterview {
  id: string
  applicationId: string
  roundName: string
  scheduledAt: string | null
  durationMinutes: number | null
  format: string
  locationOrLink: string
  result: string
  notes: string
  createdAt: string
  updatedAt: string
}

export interface InterviewRetrospective {
  interviewId: string
  questions: string[]
  selfAssessment: string
  strengths: string[]
  weaknesses: string[]
  followUpActions: string[]
  aiAnalysis: string
  createdAt: string
  updatedAt: string
}

export interface Offer {
  applicationId: string
  receivedAt: string
  status: 'pending' | 'accepted' | 'declined' | 'expired'
  deadline: string | null
  salarySummary: string
  notes: string
  decidedAt: string | null
  createdAt: string
  updatedAt: string
}

export interface CareerWorkspace {
  resumes: Resume[]
  jobs: Job[]
  applications: Application[]
  events: ApplicationEvent[]
  interviews: RealInterview[]
  retrospectives: InterviewRetrospective[]
  offers: Offer[]
}

export interface CreateResumeRequest {
  title: string
  content: string
  skills: string[]
  projectSummary: string
  note: string
  isDefault: boolean
}

export interface CreateJobRequest {
  companyName: string
  jobTitle: string
  location: string
  sourceUrl: string
  jdContent: string
}
