export function useApi() {
  const config = useRuntimeConfig()
  const auth = useAuthStore()

  async function api<T>(path: string, opts: Record<string, any> = {}): Promise<T> {
    const { _retried, headers, ...rest } = opts
    const run = () =>
      $fetch<T>(path, {
        baseURL: config.public.apiBase,
        ...rest,
        headers: {
          ...(headers || {}),
          ...(auth.accessToken ? { Authorization: `Bearer ${auth.accessToken}` } : {}),
        },
      })

    try {
      return await run()
    } catch (e: any) {
      const status = e?.response?.status ?? e?.statusCode
      if (status === 401 && auth.refreshToken && !_retried) {
        const ok = await auth.refresh()
        if (ok) return api<T>(path, { ...opts, _retried: true })
        if (import.meta.client) await navigateTo('/login')
      }
      throw e
    }
  }

  return api
}
