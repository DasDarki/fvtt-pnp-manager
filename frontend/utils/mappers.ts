import type { Character, ImageEntry, Item, Memory, MemoryLevel, SceneSummary } from '~/types/entities'
import type { ApiCharacter, ApiImage, ApiItem, ApiMemory, ApiScene } from '~/types/api'

const RINGS = [
  'conic-gradient(from 140deg,var(--primary),var(--secondary),var(--magenta),var(--primary))',
  'conic-gradient(from 120deg,var(--primary),var(--emerald),var(--primary))',
  'conic-gradient(from 90deg,var(--secondary),var(--magenta),var(--secondary))',
  'conic-gradient(from 60deg,var(--gold),var(--ember),var(--gold))',
  'conic-gradient(from 200deg,var(--primary),var(--secondary),var(--primary))',
]

function ringFor(id: string): string {
  let h = 0
  for (const ch of id) h = (h * 31 + ch.charCodeAt(0)) >>> 0
  return RINGS[h % RINGS.length]
}

export function toCharacterVM(a: ApiCharacter): Character {
  const d = a.systemData || {}
  return {
    id: a.id,
    name: a.name,
    subtitle: d.subtitle || a.summary || '',
    type: (a.characterType as Character['type']) || 'npc',
    status: (a.status as Character['status']) || 'alive',
    initial: (a.name || '?').charAt(0).toUpperCase(),
    ring: d.ring || ringFor(a.id),
    image: a.imageUrl || undefined,
    imageAlign: d.imageAlign || 'center',
    stats: Array.isArray(d.stats) ? d.stats : [],
    hpPercent: typeof d.hpPercent === 'number' ? d.hpPercent : 100,
    critical: !!d.critical,
  }
}

export function toItemVM(a: ApiItem): Item {
  const d = a.systemData || {}
  return {
    id: a.id,
    name: a.name,
    type: a.itemType || d.type || '',
    rarity: (a.rarity as Item['rarity']) || 'common',
    icon: d.icon || 'lucide:gem',
    attuned: a.attuned,
    note: a.summary,
    image: a.imageUrl || undefined,
    imageAlign: d.imageAlign || 'center',
  }
}

export function toSceneVM(a: ApiScene): SceneSummary {
  const d = a.systemData || {}
  return {
    id: a.id,
    title: a.name,
    description: a.summary || d.description || '',
    act: d.act || '',
    status: d.status || a.sceneStatus || '',
    actors: Array.isArray(d.actors) ? d.actors : [],
    extra: typeof d.extra === 'number' ? d.extra : 0,
    tone: d.tone || 'arcane',
    image: a.imageUrl || undefined,
    imageAlign: d.imageAlign || 'center',
  }
}

export function toImageVM(a: ApiImage): ImageEntry {
  return {
    id: a.id,
    name: a.name,
    image: a.imageUrl || undefined,
    imageAlign: a.imageAlign || 'center',
    pushAs: a.pushAs || 'empty_actor',
    notes: a.notes,
  }
}

export function toMemoryVM(a: ApiMemory): Memory {
  return {
    id: a.id,
    title: a.title,
    body: a.body,
    level: (a.level as MemoryLevel) || 'info',
    subjectType: (a.subjectType as Memory['subjectType']) || 'campaign',
    subjectLabel: a.subjectLabel || '',
    time: relativeTime(a.createdAt),
    acknowledged: a.acknowledged,
    pinned: a.pinned,
  }
}

export function relativeTime(iso: string): string {
  const then = new Date(iso).getTime()
  if (Number.isNaN(then)) return ''
  const diff = Date.now() - then
  const min = Math.floor(diff / 60000)
  if (min < 1) return 'gerade eben'
  if (min < 60) return `vor ${min} Min`
  const h = Math.floor(min / 60)
  if (h < 24) return `vor ${h} Std`
  const d = Math.floor(h / 24)
  return `vor ${d} Tg`
}
