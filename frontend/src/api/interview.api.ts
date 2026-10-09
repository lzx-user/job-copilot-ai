import type { ApiResponse } from '../types/api'
import type {
  InterviewOptionsResult,
  InterviewHistoryResult,
  InterviewReportResult,
  InterviewSessionDetail,
  InterviewTurnRequest,
  InterviewTurnResult,
  StartInterviewResult,
	InterviewType,
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

export async function submitInterviewTurn(payload: InterviewTurnRequest): Promise<InterviewTurnResult> {
  const response = await request.post<ApiResponse<InterviewTurnResult>>(
    '/ai/interview/turn',
    payload,
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

export async function startInterview(analysisId: string, interviewType: InterviewType, applicationId = '', resumeVersionId = ''): Promise<StartInterviewResult> {
  const response = await request.post<ApiResponse<StartInterviewResult>>(
    '/ai/interview/start',
		{ analysisId, interviewType, applicationId, resumeVersionId },
    { timeout: 100_000 },
  )
  return response.data.data
}

export async function listInterviewSessions(): Promise<InterviewHistoryResult> {
  const response = await request.get<ApiResponse<InterviewHistoryResult>>('/ai/interview/sessions')
  return response.data.data
}

export async function deleteInterviewSession(sessionId: string): Promise<void> {
  await request.delete(`/ai/interview/sessions/${sessionId}`)
}
