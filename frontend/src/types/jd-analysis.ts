export interface AnalyzeJDRequest {
  companyName: string
  jobTitle: string
  jdContent: string
  resumeSummary: string
  skills: string[]
}

export interface JdAnalysisResult {
  analysisId: string
  matchScore: number
  jobSummary: string
  coreRequirements: string[]
  matchedSkills: string[]
  missingSkills: string[]
  resumeSuggestions: string[]
  preparationTopics: string[]
  greetingMessage: string
}

export interface JdAnalysisRecord extends JdAnalysisResult {
  companyName: string
  jobTitle: string
  createdAt: string
}

export interface JdAnalysisHistoryResult {
  items: JdAnalysisRecord[]
}
