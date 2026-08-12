import { Router } from 'express'
import type { HealthResponse } from '../types/api.js'

export const healthRouter = Router()

healthRouter.get('/', (_request, response) => {
  const body: HealthResponse = {
    status: 'ok',
    message: 'Job Copilot API is running',
    timestamp: new Date().toISOString(),
  }

  response.status(200).json(body)
})

