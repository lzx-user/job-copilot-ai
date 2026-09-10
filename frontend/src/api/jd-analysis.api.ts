import type { ApiResponse } from '../types/api'
import type {
  AnalyzeJDRequest,
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
