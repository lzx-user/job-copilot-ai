import { AxiosError } from 'axios'

const AUTH_ERROR_MESSAGES: Array<[RegExp, string]> = [
  [/invalid login credentials/i, '邮箱或密码错误，请检查后重试'],
  [/email not confirmed/i, '邮箱尚未验证，请先前往邮箱完成验证'],
  [/user already registered|already been registered/i, '该邮箱已经注册，请直接登录'],
  [/password.*(?:short|least)/i, '密码长度不足，请至少输入 8 位'],
  [/signup is disabled/i, '当前 Supabase 项目未开放邮箱注册'],
  [/rate limit/i, '操作过于频繁，请稍后再试'],
]

/** 把第三方英文错误转换成用户能理解、且不会泄露服务端细节的中文提示。 */
export function toUserMessage(error: unknown): string {
  if (error instanceof AxiosError) {
    if (error.code === 'ECONNABORTED') return '请求超时，请稍后重试'
    if (!error.response) return '网络连接异常，请检查后重试'
    return '后端服务暂时无法访问，请稍后重试'
  }

  const message = error instanceof Error ? error.message : String(error ?? '')

  if (/supabase.*(?:not configured|未配置)/i.test(message)) {
    return 'Supabase 尚未配置，请先在 frontend/.env.local 中填写项目地址和 anon key'
  }

  const matched = AUTH_ERROR_MESSAGES.find(([pattern]) => pattern.test(message))
  return matched?.[1] ?? '操作失败，请稍后重试'
}

