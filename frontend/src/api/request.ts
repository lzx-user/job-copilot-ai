import axios from 'axios'
import { supabase } from '../lib/supabase'

// Axios 实例集中管理后端地址和超时，页面不需要重复拼接 URL。
export const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10_000,
  headers: {
    'Content-Type': 'application/json',
  },
})

request.interceptors.request.use(async (config) => {
  if (!supabase) return config

  const { data, error } = await supabase.auth.getSession()
  if (error) throw error
  if (data.session?.access_token) {
    config.headers.Authorization = `Bearer ${data.session.access_token}`
  }
  return config
})
