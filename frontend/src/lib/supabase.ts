import { createClient, type SupabaseClient } from '@supabase/supabase-js'

// Vite 只会把 VITE_ 前缀的变量暴露给浏览器，所以这里显式读取这两个公开配置。
const supabaseUrl = import.meta.env.VITE_SUPABASE_URL?.trim()
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY?.trim()

export const isSupabaseConfigured = Boolean(supabaseUrl && supabaseAnonKey)

/*
 * anon key 是浏览器端的公开项目标识，能否读取数据仍由登录身份和 RLS 决定。
 * 它不等于拥有管理员权限的 service role key；后者会绕过 RLS，绝不能放进前端。
 * 因此后续创建业务表时，必须为每张表开启 RLS 并编写“只访问本人数据”的策略。
 */
export const supabase: SupabaseClient | null = isSupabaseConfigured
  ? createClient(supabaseUrl as string, supabaseAnonKey as string, {
      auth: {
        persistSession: true,
        autoRefreshToken: true,
        detectSessionInUrl: true,
      },
    })
  : null

