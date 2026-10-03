import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { HttpsProxyAgent } from 'https-proxy-agent'

export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, '.', '')
  const supabaseUrl = env.VITE_SUPABASE_URL?.trim()
  const supabaseProxyUrl = env.SUPABASE_PROXY_URL?.trim()

  // 构建产物不能悄悄回退到访问者的 localhost；开发服务器仍可使用默认地址。
  if (command === 'build' && !env.VITE_API_BASE_URL?.trim()) {
    throw new Error('构建前请设置 VITE_API_BASE_URL：本地预览可用 http://localhost:8080/api/v1，生产部署请使用正式后端 HTTPS 地址。')
  }

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
              // Vite 的反向代理不会自动使用系统网络代理，需显式指定出站 Agent。
              agent: supabaseProxyUrl ? new HttpsProxyAgent(supabaseProxyUrl) : undefined,
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
