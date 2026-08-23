import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { isSupabaseConfigured, supabase } from '../lib/supabase'
import type {
  CandidateProfile,
  ProfileRow,
  SaveProfileResult,
} from '../types/profile'
import { toUserMessage } from '../utils/error'
import { useAuthStore } from './auth.store'

function createEmptyProfile(): CandidateProfile {
  return {
    nickname: '',
    targetRoles: [],
    expectedCities: [],
    skills: [],
    projectSummary: '',
    strengths: '',
    availability: '',
    graduationYear: '',
  }
}

// 数据库使用 snake_case，页面使用 camelCase，因此在数据边界显式转换。
function toCandidateProfile(row: ProfileRow): CandidateProfile {
  return {
    nickname: row.nickname,
    targetRoles: row.target_roles,
    expectedCities: row.expected_cities,
    skills: row.skills,
    projectSummary: row.project_summary,
    strengths: row.strengths,
    availability: row.availability,
    graduationYear: row.graduation_year ?? '',
  }
}

export const useProfileStore = defineStore('profile', () => {
  let authVersion = 0
  const authStore = useAuthStore()

  const profile = ref<CandidateProfile>(createEmptyProfile())
  const loading = ref(false)
  const saving = ref(false)
  const loaded = ref(false)
  const profileError = ref('')

  // 清空上一位用户的 Profile 和请求状态
  function resetProfileState() {
    profile.value = createEmptyProfile()
    loading.value = false
    saving.value = false
    loaded.value = false
    profileError.value = ''
  }

  // 判断异步请求是否仍属于当前认证周期
  function isCurrentUserRequest(
    requestUserId: string,
    requestAuthVersion: number,
  ): boolean {
    return (
      authStore.user?.id === requestUserId && authVersion === requestAuthVersion
    )
  }
  // 登录身份发生变化时重置 Profile Store
  watch(
    () => authStore.user?.id,
    (currentUserId, previousUserId) => {
      if (currentUserId === previousUserId) return
      authVersion += 1
      resetProfileState()
    },
  )

  // 从数据库读取当前登录用户的 Profile
  async function fetchProfile() {
    if (loading.value) return

    profileError.value = ''
    loaded.value = false

    if (!isSupabaseConfigured || !supabase) {
      profileError.value = toUserMessage(new Error('Supabase 未配置'))
      return
    }

    // userId 必须来自认证状态，不能来自页面表单。
    const userId = authStore.user?.id
    if (!userId) {
      profileError.value = '登录状态无效，请重新登录'
      return
    }

    const requestAuthVersion = authVersion

    loading.value = true

    try {
      const { data, error } = await supabase
        .from('profiles')
        .select(
          `
          user_id,
          nickname,
          target_roles,
          expected_cities,
          skills,
          project_summary,
          strengths,
          availability,
          graduation_year,
          created_at,
          updated_at
        `,
        )
        // 只查询当前用户对应的行；真正的安全边界仍然是 RLS
        .eq('user_id', userId)
        // 允许结果为零条或一条。新用户没有 Profile 时返回 null，不会被当成错误
        .maybeSingle<ProfileRow>()

      if (!isCurrentUserRequest(userId, requestAuthVersion)) return

      if (error) throw error

      // 新用户没有 Profile 是正常状态，应显示空表单，而不是显示系统错误。
      profile.value = data ? toCandidateProfile(data) : createEmptyProfile()
      loaded.value = true
    } catch (error) {
      if (isCurrentUserRequest(userId, requestAuthVersion)) {
        profileError.value = toUserMessage(error)
      }
    } finally {
      if (isCurrentUserRequest(userId, requestAuthVersion)) {
        loading.value = false
      }
    }
  }

  // 将当前 Profile 保存到数据库
  async function saveProfile(
    nextProfile: CandidateProfile,
  ): Promise<SaveProfileResult> {
    // 防止用户连续点击导致同一份档案被重复提交
    if (saving.value) return 'ignored'

    profileError.value = ''

    if (!isSupabaseConfigured || !supabase) {
      profileError.value = toUserMessage(new Error('Supabase 未配置'))
      return 'error'
    }

    // 数据所有者只能来自认证系统，不能由页面表单传入
    const userId = authStore.user?.id
    if (!userId) {
      profileError.value = '登录状态无效，请重新登录'
      return 'error'
    }

    const requestAuthVersion = authVersion

    saving.value = true

    try {
      const { data, error } = await supabase
        .from('profiles')
        .upsert(
          {
            user_id: userId,
            nickname: nextProfile.nickname.trim(),
            target_roles: nextProfile.targetRoles,
            expected_cities: nextProfile.expectedCities,
            skills: nextProfile.skills,
            project_summary: nextProfile.projectSummary.trim(),
            strengths: nextProfile.strengths.trim(),
            availability: nextProfile.availability.trim(),
            graduation_year: nextProfile.graduationYear.trim() || null,
          },
          {
            onConflict: 'user_id',
          },
        )
        .select(
          `
        user_id,
        nickname,
        target_roles,
        expected_cities,
        skills,
        project_summary,
        strengths,
        availability,
        graduation_year,
        created_at,
        updated_at
      `,
        )
        .single<ProfileRow>()

      if (!isCurrentUserRequest(userId, requestAuthVersion)) return 'ignored'

      if (error) throw error

      // 使用数据库返回的最终结果更新 Store，避免页面状态与数据库不一致
      profile.value = toCandidateProfile(data)
      loaded.value = true
      return 'success'
    } catch (error) {
      if (!isCurrentUserRequest(userId, requestAuthVersion)) return 'ignored'
      profileError.value = toUserMessage(error)
      return 'error'
    } finally {
      if (isCurrentUserRequest(userId, requestAuthVersion)) {
        saving.value = false
      }
    }
  }

  return {
    profile,
    loading,
    saving,
    loaded,
    profileError,
    fetchProfile,
    saveProfile,
  }
})
