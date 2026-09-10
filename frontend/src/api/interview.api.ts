import type { ApiResponse } from '../types/api'
import type { InterviewOptionsResult, StartInterviewResult } from '../types/interview'
import { request } from './request'

export async function listInterviewOptions(): Promise<InterviewOptionsResult> {
  const response = await request.get<ApiResponse<InterviewOptionsResult>>('/ai/interview/options')
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
