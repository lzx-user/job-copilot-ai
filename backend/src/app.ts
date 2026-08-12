import cors from 'cors'
import express from 'express'
import rateLimit from 'express-rate-limit'
import helmet from 'helmet'
import { env } from './config/env.js'
import { errorMiddleware } from './middleware/error.middleware.js'
import { notFoundMiddleware } from './middleware/not-found.middleware.js'
import { apiRouter } from './routes/index.js'
import { healthRouter } from './routes/health.routes.js'

export const app = express()

// 安全相关中间件放在业务路由之前，确保所有接口都得到相同保护。
app.use(helmet())
app.use(
  cors({
    origin: env.FRONTEND_ORIGIN,
    methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'],
  }),
)
app.use(express.json({ limit: '1mb' }))
app.use(
  rateLimit({
    windowMs: 15 * 60 * 1000,
    limit: 300,
    standardHeaders: 'draft-8',
    legacyHeaders: false,
    message: {
      code: 'RATE_LIMITED',
      message: '请求过于频繁，请稍后再试',
      data: null,
    },
  }),
)

app.use('/health', healthRouter)
app.use('/api/v1', apiRouter)
app.use(notFoundMiddleware)
app.use(errorMiddleware)

