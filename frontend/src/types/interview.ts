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
