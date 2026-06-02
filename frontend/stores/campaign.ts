import { defineStore } from 'pinia'
import type { ApiCampaign, ApiCharacter } from '~/types/api'

export const useCampaignStore = defineStore('campaign', () => {
  const STORAGE = 'aw-campaign'
  const current = ref<ApiCampaign | null>(null)
  const campaigns = ref<ApiCampaign[]>([])
  const ready = ref(false)
  const currentId = computed(() => current.value?.id ?? null)
  let inflight: Promise<void> | null = null

  function persist() {
    if (import.meta.client && current.value) localStorage.setItem(STORAGE, current.value.id)
  }

  function ensure(): Promise<void> {
    if (ready.value) return Promise.resolve()
    if (inflight) return inflight
    inflight = (async () => {
      const api = useApi()
      campaigns.value = await api<ApiCampaign[]>('/campaigns')
      if (!campaigns.value.length) {
        const cam = await api<ApiCampaign>('/campaigns', {
          method: 'POST',
          body: { name: 'Asche & Aether', ruleset: 'dnd5e_2024', forgeRootPath: 'worlds/aa/aetherwright' },
        })
        campaigns.value = [cam]
        current.value = cam
        await seed(cam.id)
      } else {
        const savedId = import.meta.client ? localStorage.getItem(STORAGE) : null
        current.value = campaigns.value.find((c) => c.id === savedId) || campaigns.value[0]
      }
      persist()
      ready.value = true
    })().finally(() => {
      inflight = null
    })
    return inflight
  }

  async function refreshList() {
    const api = useApi()
    campaigns.value = await api<ApiCampaign[]>('/campaigns')
    if (current.value) {
      current.value = campaigns.value.find((c) => c.id === current.value!.id) || campaigns.value[0] || null
    }
    persist()
  }

  function select(idToSelect: string) {
    const cam = campaigns.value.find((c) => c.id === idToSelect)
    if (cam) {
      current.value = cam
      persist()
    }
  }

  async function createCampaign(name: string, ruleset: string, forgeRootPath = '') {
    const api = useApi()
    const cam = await api<ApiCampaign>('/campaigns', {
      method: 'POST',
      body: { name, ruleset, forgeRootPath },
    })
    campaigns.value = [cam, ...campaigns.value]
    return cam
  }

  async function seed(cid: string) {
    const api = useApi()
    const demo = useDemoCampaign()
    const nameToId: Record<string, string> = {}
    const num = (v: any) => {
      const x = parseInt(String(v ?? '').replace(/[^0-9-]/g, ''), 10)
      return Number.isFinite(x) ? x : 0
    }

    for (const c of demo.characters) {
      const parts = c.subtitle.split(' · ')
      const hp = num(c.stats[1]?.value)
      const created = await api<ApiCharacter>(`/campaigns/${cid}/characters`, {
        method: 'POST',
        body: {
          name: c.name,
          characterType: c.type,
          status: c.status,
          systemData: {
            subtitle: c.subtitle,
            ring: c.ring,
            stats: c.stats,
            hpPercent: c.hpPercent,
            critical: c.critical || false,
            race: parts[0] || '',
            class: parts[1] || '',
            ac: num(c.stats[0]?.value),
            hp,
            hpMax: hp || 10,
            level: num(c.stats[2]?.value),
            proficiency: 3,
            abilities: { str: 10, dex: 14, con: 12, int: 13, wis: 11, cha: 16 },
          },
        },
      })
      nameToId[c.name] = created.id
    }

    for (const it of demo.items) {
      await api(`/campaigns/${cid}/items`, {
        method: 'POST',
        body: {
          name: it.name,
          itemType: it.type,
          rarity: it.rarity,
          attuned: it.attuned || false,
          summary: it.note || '',
          systemData: { icon: it.icon },
        },
      })
    }

    for (const s of demo.scenes) {
      await api(`/campaigns/${cid}/scenes`, {
        method: 'POST',
        body: {
          name: s.title,
          summary: s.description,
          sceneStatus: 'prepared',
          systemData: { act: s.act, status: s.status, actors: s.actors, extra: s.extra, tone: s.tone },
        },
      })
    }

    const demoImages = [
      { name: 'Wappen des Hauses Vael', notes: 'Silberner Rabe auf schwarzem Grund.', pushAs: 'empty_actor' },
      { name: 'Karte der Hauptstadt', notes: 'Handout für die Spieler — markierte Treffpunkte.', pushAs: 'journal' },
      { name: 'Das versiegelte Tor', notes: 'Runen, die im Mondlicht aufleuchten.', pushAs: 'journal' },
    ]
    for (const img of demoImages) {
      await api(`/campaigns/${cid}/images`, { method: 'POST', body: img })
    }

    for (const m of demo.memories) {
      const created = await api<{ id: string }>(`/campaigns/${cid}/memories`, {
        method: 'POST',
        body: {
          title: m.title,
          body: m.body,
          level: m.level,
          subjectType: m.subjectType,
          subjectLabel: m.subjectLabel,
          subjectId: nameToId[m.subjectLabel] || null,
        },
      })
      if (m.acknowledged) {
        await api(`/campaigns/${cid}/memories/${created.id}`, {
          method: 'PATCH',
          body: { acknowledged: true },
        })
      }
    }

    await api(`/campaigns/${cid}/memories`, {
      method: 'POST',
      body: {
        title: demo.criticalAlert.title,
        body: demo.criticalAlert.body,
        level: 'critical',
        subjectType: 'character',
        subjectLabel: demo.criticalAlert.subjectLabel,
        subjectId: nameToId[demo.criticalAlert.subjectLabel] || null,
      },
    })
  }

  async function updateStyle(artStyle: string, stylePrompt: string) {
    if (!currentId.value) return
    const api = useApi()
    const cam = await api<ApiCampaign>(`/campaigns/${currentId.value}/style`, {
      method: 'PATCH',
      body: { artStyle, stylePrompt },
    })
    current.value = cam
    const i = campaigns.value.findIndex((c) => c.id === cam.id)
    if (i >= 0) campaigns.value[i] = cam
  }

  function reset() {
    current.value = null
    campaigns.value = []
    ready.value = false
  }

  return {
    current,
    campaigns,
    currentId,
    ready,
    ensure,
    refreshList,
    select,
    createCampaign,
    updateStyle,
    reset,
  }
})
