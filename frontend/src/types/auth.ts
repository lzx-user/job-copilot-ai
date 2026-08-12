import type { Session, User } from '@supabase/supabase-js'

export type AuthSession = Session
export type AuthUser = User

export interface AuthActionResult {
  success: boolean
  requiresEmailConfirmation?: boolean
}

export interface LoginFormValues {
  email: string
  password: string
}

export interface RegisterFormValues extends LoginFormValues {
  confirmPassword: string
  acceptedTerms: boolean
}

