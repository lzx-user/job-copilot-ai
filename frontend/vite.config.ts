import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    // 与 backend/.env.example 的 FRONTEND_ORIGIN 保持一致，避免本地开发时触发 CORS。
    host: 'localhost',
    port: 5173,
  },
})
