import { defineStore } from 'pinia'
import type { ApiUser, AuthSession } from '~/types/api'

const STORAGE_KEY = 'aw-auth'

export const useAuthStore = defineStore('auth', () => {
  const config = useRuntimeConfig()
  const base = config.public.apiBase

  const user = ref<ApiUser | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)

  const isAuthed = computed(() => !!accessToken.value)

  function persist() {
    if (!import.meta.client) return
    if (accessToken.value && refreshToken.value) {
      localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({ a: accessToken.value, r: refreshToken.value, u: user.value }),
      )
    } else {
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  function hydrate() {
    if (!import.meta.client) return
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return
    try {
      const d = JSON.parse(raw)
      accessToken.value = d.a
      refreshToken.value = d.r
      user.value = d.u
    } catch {
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  function setSession(s: AuthSession) {
    user.value = s.user
    accessToken.value = s.accessToken
    refreshToken.value = s.refreshToken
    persist()
  }

  function clear() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    persist()
  }

  function device() {
    if (import.meta.client) return navigator.userAgent.slice(0, 80)
    return 'server'
  }

  async function login(email: string, password: string) {
    const s = await $fetch<AuthSession>(`${base}/auth/login`, {
      method: 'POST',
      body: { email, password, device: device() },
    })
    setSession(s)
  }

  async function register(email: string, password: string, name: string) {
    const s = await $fetch<AuthSession>(`${base}/auth/register`, {
      method: 'POST',
      body: { email, password, name, device: device() },
    })
    setSession(s)
  }

  let refreshInflight: Promise<boolean> | null = null
  async function refresh(): Promise<boolean> {
    // Dedupe concurrent refreshes: the backend rotates (revokes) the refresh
    // token, so parallel 401s must share a single refresh call — otherwise the
    // second one presents an already-revoked token and forces a spurious logout.
    if (refreshInflight) return refreshInflight
    refreshInflight = (async () => {
      if (!refreshToken.value) return false
      try {
        const s = await $fetch<AuthSession>(`${base}/auth/refresh`, {
          method: 'POST',
          body: { refreshToken: refreshToken.value, device: device() },
        })
        setSession(s)
        return true
      } catch {
        clear()
        return false
      }
    })().finally(() => {
      refreshInflight = null
    })
    return refreshInflight
  }

  async function logout() {
    const token = refreshToken.value
    clear()
    if (token) {
      try {
        await $fetch(`${base}/auth/logout`, { method: 'POST', body: { refreshToken: token } })
      } catch {
        // ignore
      }
    }
  }

  return {
    user,
    accessToken,
    refreshToken,
    isAuthed,
    hydrate,
    setSession,
    clear,
    login,
    register,
    refresh,
    logout,
  }
})
