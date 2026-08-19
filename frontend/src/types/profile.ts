// 页面和 Store 使用的 camelCase 业务数据，不包含可编辑的 userId
export interface CandidateProfile {
  nickname: string
  targetRoles: string[]
  expectedCities: string[]
  skills: string[]
  projectSummary: string
  strengths: string
  availability: string
  graduationYear: string
}

// Supabase 数据库返回的真实 snake_case 行结构，包含所有者和时间字段
export interface ProfileRow {
  user_id: string
  nickname: string
  target_roles: string[]
  expected_cities: string[]
  skills: string[]
  project_summary: string
  strengths: string
  availability: string
  graduation_year: string | null
  created_at: string
  updated_at: string
}

// 结果类型，表示保存档案的结果状态
export type SaveProfileResult = 'success' | 'error' | 'ignored'
