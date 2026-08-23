export interface ApiResponse<T> {
  code: string
  message?: string
  data: T
}

export interface HealthData {
  status: 'ok'
}

export type ApiErrorResponse = ApiResponse<null>

export type ServiceStatus = 'checking' | 'online' | 'offline'
