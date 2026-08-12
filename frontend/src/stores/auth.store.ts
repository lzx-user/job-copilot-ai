import type { Session, User } from '@supabase/supabase-js'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { isSupabaseConfigured, supabase } from '../lib/supabase'
import type { AuthActionResult } from '../types/auth'
import { toUserMessage } from '../utils/error'

let removeAuthListener: (() => void) | null = null

export const useAuthStore = defineStore('auth', () => {
  // Pinia 的 state 保存跨页面共享的登录状态；action 负责有副作用的认证请求。
  const session = ref<Session | null>(null)
  const user = ref<User | null>(null)
  const initializing = ref(true)
  const submitting = ref(false)
  const initialized = ref(false)
  const authError = ref('')

  const isAuthenticated = computed(() => Boolean(session.value))

  function clearError() {
    authError.value = ''
  }

  async function initialize() {
    if (initialized.value) return

    initializing.value = true
    clearError()

    try {
      if (!isSupabaseConfigured || !supabase) {
        session.value = null
        user.value = null
        return
      }

      // Session 包含登录有效期等认证信息，User 是其中的用户资料；路由判断应以 Session 为准。
      const { data, error } = await supabase.auth.getSession()
      if (error) throw error
      session.value = data.session
      user.value = data.session?.user ?? null

      // onAuthStateChange 会同步其他标签页中的登录/退出。先取消旧监听，避免热更新时重复触发。
      removeAuthListener?.()
      const { data: listener } = supabase.auth.onAuthStateChange((_event, nextSession) => {
        session.value = nextSession
        user.value = nextSession?.user ?? null
      })
      removeAuthListener = () => listener.subscription.unsubscribe()
    } catch (error) {
      authError.value = toUserMessage(error)
      session.value = null
      user.value = null
    } finally {
      // 初始化完成前路由不能贸然跳转，否则刷新时会误把已登录用户送回登录页。
      initialized.value = true
      initializing.value = false
    }
  }

  async function login(email: string, password: string): Promise<AuthActionResult> {
    if (submitting.value) return { success: false }
    clearError()

    if (!isSupabaseConfigured || !supabase) {
      authError.value = toUserMessage(new Error('Supabase 未配置'))
      return { success: false }
    }

    submitting.value = true
    try {
      const { data, error } = await supabase.auth.signInWithPassword({ email, password })
      if (error) throw error
      session.value = data.session
      user.value = data.user
      return { success: true }
    } catch (error) {
      authError.value = toUserMessage(error)
      return { success: false }
    } finally {
      submitting.value = false
    }
  }

  async function register(email: string, password: string): Promise<AuthActionResult> {
    if (submitting.value) return { success: false }
    clearError()

    if (!isSupabaseConfigured || !supabase) {
      authError.value = toUserMessage(new Error('Supabase 未配置'))
      return { success: false }
    }

    submitting.value = true
    try {
      const { data, error } = await supabase.auth.signUp({ email, password })
      if (error) throw error
      session.value = data.session
      user.value = data.user
      return {
        success: true,
        requiresEmailConfirmation: !data.session,
      }
    } catch (error) {
      authError.value = toUserMessage(error)
      return { success: false }
    } finally {
      submitting.value = false
    }
  }

  async function logout(): Promise<AuthActionResult> {
    if (submitting.value) return { success: false }
    clearError()

    if (!supabase) {
      session.value = null
      user.value = null
      return { success: true }
    }

    submitting.value = true
    try {
      const { error } = await supabase.auth.signOut()
      if (error) throw error
      session.value = null
      user.value = null
      return { success: true }
    } catch (error) {
      authError.value = toUserMessage(error)
      return { success: false }
    } finally {
      submitting.value = false
    }
  }

  function disposeAuthListener() {
    removeAuthListener?.()
    removeAuthListener = null
  }

  return {
    session,
    user,
    initializing,
    submitting,
    initialized,
    authError,
    isAuthenticated,
    initialize,
    login,
    register,
    logout,
    clearError,
    disposeAuthListener,
  }
})

