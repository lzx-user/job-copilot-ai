import type { ApiResponse } from '../types/api'
import type { DashboardSummary } from '../types/dashboard'
import { request } from './request'

export async function getDashboardSummary(): Promise<DashboardSummary> {
  const response = await request.get<ApiResponse<DashboardSummary>>('/ai/dashboard')
  return response.data.data
}
