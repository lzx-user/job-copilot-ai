export interface ApiErrorResponse {
  code: string
  message: string
  data: null
}

export interface HealthResponse {
  status: 'ok'
  message: string
  timestamp: string
}

export type ServiceStatus = 'checking' | 'online' | 'offline'

