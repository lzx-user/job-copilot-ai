import axios from 'axios'

// Axios 实例集中管理后端地址和超时，页面不需要重复拼接 URL。
export const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api/v1',
  timeout: 10_000,
  headers: {
    'Content-Type': 'application/json',
  },
})

