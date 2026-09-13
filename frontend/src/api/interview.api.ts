import type { ApiResponse } from '../types/api'
import type {
  InterviewOptionsResult,
  InterviewReportResult,
  InterviewSessionDetail,
  InterviewTurnResult,
  StartInterviewResult,
} from '../types/interview'
import { request } from './request'

export async function listInterviewOptions(): Promise<InterviewOptionsResult> {
  const response = await request.get<ApiResponse<InterviewOptionsResult>>('/ai/interview/options')
  return response.data.data
}

export async function getInterviewSession(sessionId: string): Promise<InterviewSessionDetail> {
  const response = await request.get<ApiResponse<InterviewSessionDetail>>(`/ai/interview/sessions/${sessionId}`)
  return response.data.data
}

export async function submitInterviewTurn(sessionId: string, answer: string): Promise<InterviewTurnResult> {
  const response = await request.post<ApiResponse<InterviewTurnResult>>(
    '/ai/interview/turn',
    { sessionId, answer },
    { timeout: 100_000 },
  )
  return response.data.data
}

export async function generateInterviewReport(sessionId: string): Promise<InterviewReportResult> {
  const response = await request.post<ApiResponse<InterviewReportResult>>(
    '/ai/interview/report',
    { sessionId },
    { timeout: 100_000 },
  )
  return response.data.data
}

export async function startInterview(analysisId: string): Promise<StartInterviewResult> {
  const response = await request.post<ApiResponse<StartInterviewResult>>(
    '/ai/interview/start',
    { analysisId },
    { timeout: 100_000 },
  )
  return response.data.data
}
