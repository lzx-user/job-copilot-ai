import type { ApiResponse } from '../types/api'
import type { CareerWorkspace, CreateJobRequest, CreateResumeRequest } from '../types/career'
import { request } from './request'

export async function getCareerWorkspace(): Promise<CareerWorkspace> {
  const response = await request.get<ApiResponse<CareerWorkspace>>('/career/workspace')
  return response.data.data
}

export async function createResume(payload: CreateResumeRequest) {
  const response = await request.post<ApiResponse<{ resumeId: string; versionId: string }>>('/career/resumes', payload)
  return response.data.data
}

export async function createResumeVersion(resumeId: string, payload: Omit<CreateResumeRequest, 'title' | 'isDefault'>) {
  const response = await request.post<ApiResponse<{ versionId: string }>>(`/career/resumes/${resumeId}/versions`, payload)
  return response.data.data
}

export async function createJob(payload: CreateJobRequest) {
  const response = await request.post<ApiResponse<{ jobId: string; jdVersionId: string }>>('/career/jobs', payload)
  return response.data.data
}

export async function createJDVersion(jobId: string, jdContent: string) {
  const response = await request.post<ApiResponse<{ jdVersionId: string }>>(`/career/jobs/${jobId}/jd-versions`, { jdContent })
  return response.data.data
}

export async function createApplication(payload: Record<string, unknown>) {
  const response = await request.post<ApiResponse<{ applicationId: string }>>('/career/applications', payload)
  return response.data.data
}

export async function addApplicationEvent(payload: Record<string, unknown>) {
  const response = await request.post<ApiResponse<{ eventId: string }>>('/career/events', payload)
  return response.data.data
}

export async function createRealInterview(payload: Record<string, unknown>) {
  const response = await request.post<ApiResponse<{ interviewId: string }>>('/career/interviews', payload)
  return response.data.data
}

export async function updateRealInterview(interviewId: string, payload: { result: string; notes: string }) {
  await request.patch(`/career/interviews/${interviewId}`, payload)
}

export async function saveRetrospective(payload: Record<string, unknown>) {
  await request.put('/career/retrospectives', payload)
}

export async function saveOffer(payload: Record<string, unknown>) {
  await request.put('/career/offers', payload)
}
