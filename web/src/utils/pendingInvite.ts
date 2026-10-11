// The invitation token a visitor arrived with, kept for the session so it survives signing in or a
// round trip through SSO. Storage can be unavailable (private mode, blocked site data); the flow
// then just loses the token and the invitee can still accept from the in-app pending list.
const KEY = 'mb_pending_invite'

export function savePendingInvite(token: string) {
  try {
    sessionStorage.setItem(KEY, token)
  } catch { /* storage unavailable */ }
}

export function pendingInvite(): string | null {
  try {
    return sessionStorage.getItem(KEY)
  } catch {
    return null
  }
}

export function clearPendingInvite() {
  try {
    sessionStorage.removeItem(KEY)
  } catch { /* storage unavailable */ }
}
