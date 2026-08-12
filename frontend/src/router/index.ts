import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth.store'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'root',
      component: () => import('../views/RootRedirectView.vue'),
    },
    {
      path: '/auth/login',
      name: 'login',
      component: () => import('../views/auth/LoginView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/auth/register',
      name: 'register',
      component: () => import('../views/auth/RegisterView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/app',
      component: () => import('../layouts/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/app/dashboard' },
        { path: 'dashboard', name: 'dashboard', component: () => import('../views/dashboard/DashboardView.vue') },
        { path: 'jd-analysis', name: 'jd-analysis', component: () => import('../views/jd-analysis/JdAnalysisView.vue') },
        { path: 'interviews', name: 'interviews', component: () => import('../views/interview/InterviewListView.vue') },
        { path: 'interviews/:id', name: 'interview-session', component: () => import('../views/interview/InterviewSessionView.vue') },
        { path: 'history', name: 'history', component: () => import('../views/history/HistoryView.vue') },
        { path: 'profile', name: 'profile', component: () => import('../views/profile/ProfileView.vue') },
        { path: 'settings', name: 'settings', component: () => import('../views/settings/SettingsView.vue') },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../views/NotFoundView.vue'),
    },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  // 每次导航先确认认证初始化完成，守卫才能根据真实 Session 做决定。
  if (!authStore.initialized) await authStore.initialize()

  if (to.path === '/') {
    return authStore.session ? '/app/dashboard' : '/auth/login'
  }

  if (to.meta.requiresAuth && !authStore.session) {
    return {
      path: '/auth/login',
      query: { redirect: to.fullPath },
    }
  }

  if (to.meta.guestOnly && authStore.session) return '/app/dashboard'

  return true
})

export default router

