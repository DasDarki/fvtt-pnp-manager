import { defineStore } from 'pinia'
import type { ApiProvider, ProvidersResponse } from '~/types/api'

export const useProviderStore = defineStore('providers', () => {
  const STORAGE = 'aw-image-provider'
  const providers = ref<ApiProvider[]>([])
  const catalog = ref<string[]>([])
  const selected = ref('')
  const ready = ref(false)
  let inflight: Promise<void> | null = null

  const available = computed(() => providers.value.filter((p) => p.hasKey))

  function reconcile() {
    if (selected.value && !available.value.some((p) => p.provider === selected.value)) {
      selected.value = ''
      if (import.meta.client) localStorage.removeItem(STORAGE)
    }
  }

  function apply(r: ProvidersResponse) {
    providers.value = r.providers
    catalog.value = r.catalog
    if (import.meta.client) {
      const saved = localStorage.getItem(STORAGE) || ''
      if (saved && r.providers.some((p) => p.provider === saved && p.hasKey)) selected.value = saved
    }
    reconcile()
  }

  function ensure(): Promise<void> {
    if (ready.value) return Promise.resolve()
    if (inflight) return inflight
    inflight = (async () => {
      const api = useApi()
      apply(await api<ProvidersResponse>('/providers'))
      ready.value = true
    })().finally(() => {
      inflight = null
    })
    return inflight
  }

  async function reload() {
    const api = useApi()
    apply(await api<ProvidersResponse>('/providers'))
  }

  function select(provider: string) {
    selected.value = provider
    if (!import.meta.client) return
    if (provider) localStorage.setItem(STORAGE, provider)
    else localStorage.removeItem(STORAGE)
  }

  async function save(provider: string, model: string, apiKey: string) {
    const api = useApi()
    const r = await api<ApiProvider>(`/providers/${provider}`, { method: 'PUT', body: { model, apiKey } })
    const i = providers.value.findIndex((p) => p.provider === provider)
    if (i >= 0) providers.value[i] = r
    else providers.value.push(r)
    return r
  }

  async function remove(provider: string) {
    const api = useApi()
    await api(`/providers/${provider}`, { method: 'DELETE' })
    providers.value = providers.value.filter((p) => p.provider !== provider)
    if (selected.value === provider) select('')
  }

  function reset() {
    providers.value = []
    catalog.value = []
    selected.value = ''
    ready.value = false
  }

  return { providers, catalog, selected, ready, available, ensure, reload, select, save, remove, reset }
})
