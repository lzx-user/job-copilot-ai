import type { ApiResponse, HealthData } from '../types/api'
import { request } from './request'

export async function getHealthStatus() {
  const response = await request.get<ApiResponse<HealthData>>('/health')
  return response.data.data
}
