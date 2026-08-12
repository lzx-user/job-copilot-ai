import type { ErrorRequestHandler } from 'express'
import type { ApiErrorResponse } from '../types/api.js'

export const errorMiddleware: ErrorRequestHandler = (error, _request, response, _next) => {
  // 仅在服务端记录必要错误。浏览器只接收稳定结构，避免泄露调用栈。
  console.error('请求处理失败：', error instanceof Error ? error.message : '未知错误')

  const body: ApiErrorResponse = {
    code: 'INTERNAL_ERROR',
    message: '服务出现异常，请稍后重试',
    data: null,
  }

  response.status(500).json(body)
}

