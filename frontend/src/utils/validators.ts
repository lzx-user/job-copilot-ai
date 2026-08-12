export const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function isValidEmail(value: string) {
  return EMAIL_PATTERN.test(value.trim())
}

export function isSafeInternalRedirect(value: unknown): value is string {
  return (
    typeof value === 'string' &&
    value.startsWith('/') &&
    !value.startsWith('//') &&
    !value.includes('\\') &&
    !/^(?:\/)?(?:https?:|javascript:)/i.test(value)
  )
}

