export interface DashboardRecentRecord {
  id: string
  kind: 'jd' | 'interview'
  companyName: string
  jobTitle: string
  description: string
  createdAt: string
}

export interface DashboardSummary {
  jdCount: number
  interviewCount: number
  applicationCount: number
  activeApplicationCount: number
  realInterviewCount: number
  offerCount: number
  applicationToInterviewRate: number
  interviewToOfferRate: number
  weeklyRecordCount: number
  profileCompleteness: number
  recentRecords: DashboardRecentRecord[]
}
