import { createPinia } from 'pinia'
import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth.store'
import './styles/variables.css'
import './styles/reset.css'
import './styles/global.css'
import './styles/element-overrides.css'

async function bootstrap() {
  const app = createApp(App)
  const pinia = createPinia()

  app.use(pinia)
  app.use(router)
  app.use(ElementPlus)

  const authStore = useAuthStore(pinia)
  // async/await 让首次 Session 检查结束后再挂载，避免页面刷新时闪过错误的登录页。
  await authStore.initialize()

  app.mount('#app')
}

void bootstrap()

