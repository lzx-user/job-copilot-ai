import type { ApiResponse } from '../types/api'
import type {
  AnalyzeJDRequest,
  JdAnalysisHistoryResult,
  JdAnalysisRecord,
  JdAnalysisResult,
} from '../types/jd-analysis'
import { request } from './request'

export async function analyzeJD(
  payload: AnalyzeJDRequest,
): Promise<JdAnalysisResult> {
  const response = await request.post<ApiResponse<JdAnalysisResult>>(
    '/ai/analyze-jd',
    payload,
    { timeout: 100_000 },
  )
  return response.data.data
}

export async function listJDAnalyses(): Promise<JdAnalysisHistoryResult> {
  const response = await request.get<ApiResponse<JdAnalysisHistoryResult>>('/ai/jd-analyses')
  return response.data.data
}

export async function getJDAnalysis(analysisId: string): Promise<JdAnalysisRecord> {
  const response = await request.get<ApiResponse<JdAnalysisRecord>>(`/ai/jd-analyses/${analysisId}`)
  return response.data.data
}
