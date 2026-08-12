import 'dotenv/config'
import { z } from 'zod'

const envSchema = z.object({
  NODE_ENV: z.enum(['development', 'test', 'production']).default('development'),
  PORT: z.coerce.number().int().positive().default(3000),
  FRONTEND_ORIGIN: z.string().url().default('http://localhost:5173'),
})

const parsedEnv = envSchema.safeParse(process.env)

if (!parsedEnv.success) {
  // 启动阶段只报告环境变量名称和校验问题，不输出任何密钥值。
  console.error('后端环境变量配置无效：', parsedEnv.error.flatten().fieldErrors)
  process.exit(1)
}

export const env = parsedEnv.data

