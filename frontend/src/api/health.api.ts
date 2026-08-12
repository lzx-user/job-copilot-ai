import type { HealthResponse } from '../types/api'
import { request } from './request'

export async function getHealthStatus() {
  const response = await request.get<HealthResponse>('/health')
  return response.data
}

