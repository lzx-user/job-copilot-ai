export interface InterviewOption {
  analysisId: string
  companyName: string
  jobTitle: string
  matchScore: number
  createdAt: string
}

export interface InterviewOptionsResult {
  items: InterviewOption[]
}

export interface StartInterviewResult {
  sessionId: string
  status: 'in_progress'
  currentRound: number
  maxRounds: number
  question: string
}

export interface InterviewFeedback {
  score: number
  feedback: string
  strengths: string[]
  improvements: string[]
}

export interface InterviewMessage extends Partial<InterviewFeedback> {
  id: string
  role: 'interviewer' | 'candidate'
  round: number
  content: string
  createdAt: string
}

export interface InterviewSessionDetail {
  sessionId: string
  status: 'pending' | 'in_progress' | 'completed'
  currentRound: number
  maxRounds: number
  companyName: string
  jobTitle: string
  messages: InterviewMessage[]
  report: InterviewReport | null
}

export interface InterviewTurnResult extends InterviewFeedback {
  sessionId: string
  currentRound: number
  maxRounds: number
  nextQuestion: string | null
}

export interface InterviewReport {
  overallScore: number
  technicalScore: number
  expressionScore: number
  projectDepthScore: number
  strengths: string[]
  weaknesses: string[]
  recommendedTopics: string[]
  answerTips: string[]
  summary: string
}

export interface InterviewReportResult extends InterviewReport {
  sessionId: string
  status: 'completed'
}
