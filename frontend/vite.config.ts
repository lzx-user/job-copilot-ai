import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, '.', '')
  const supabaseUrl = env.VITE_SUPABASE_URL?.trim()

  return {
    plugins: [vue()],
    server: {
      // 与 backend/.env.example 的 FRONTEND_ORIGIN 保持一致，避免本地开发时触发 CORS。
      host: 'localhost',
      port: 5173,
      // 本地浏览器只访问同源地址，避免代理、DNS 或扩展阻断 Supabase 域名直连。
      proxy: supabaseUrl
        ? {
            '/supabase': {
              target: supabaseUrl,
              changeOrigin: true,
              secure: true,
              ws: true,
              rewrite: (path) => path.replace(/^\/supabase/, ''),
            },
          }
        : undefined,
    },
  }
})
