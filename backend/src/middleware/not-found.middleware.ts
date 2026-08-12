import type { Request, Response } from 'express'
import type { ApiErrorResponse } from '../types/api.js'

export function notFoundMiddleware(request: Request, response: Response<ApiErrorResponse>) {
  response.status(404).json({
    code: 'NOT_FOUND',
    message: `未找到请求的接口：${request.method} ${request.path}`,
    data: null,
  })
}

