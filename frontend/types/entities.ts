export type Ruleset = 'dnd5e_2024' | 'vampire_v5'

export type CharacterStatus = 'alive' | 'dead' | 'unknown' | 'hunted'
export type CharacterType = 'pc' | 'npc' | 'ally' | 'foe' | 'neutral'

export type Rarity = 'common' | 'uncommon' | 'rare' | 'epic' | 'legendary' | 'artifact'

export type MemoryLevel = 'info' | 'notice' | 'warning' | 'critical'
export type SubjectType = 'character' | 'item' | 'scene' | 'image' | 'campaign'

export type SyncState = 'none' | 'pending' | 'synced' | 'dirty' | 'error'

export interface TagRef {
  id: string
  name: string
  color: string
}

export type AccentToken = 'primary' | 'secondary' | 'magenta' | 'gold' | 'emerald' | 'ember'

export interface MiniStat {
  label: string
  value: string
}

export interface Character {
  id: string
  name: string
  subtitle: string
  type: CharacterType
  status: CharacterStatus
  initial: string
  ring: string
  image?: string
  imageAlign?: string
  stats: MiniStat[]
  hpPercent: number
  critical?: boolean
  tags?: TagRef[]
}

export interface Memory {
  id: string
  title: string
  body: string
  level: MemoryLevel
  subjectType: SubjectType
  subjectLabel: string
  time: string
  acknowledged?: boolean
  pinned?: boolean
}

export interface Item {
  id: string
  name: string
  type: string
  rarity: Rarity
  icon: string
  attuned?: boolean
  note?: string
  image?: string
  imageAlign?: string
  tags?: TagRef[]
}

export interface ImageEntry {
  id: string
  name: string
  image?: string
  imageAlign?: string
  pushAs: string
  notes?: string
  tags?: TagRef[]
}

export type SceneTone = 'arcane' | 'ember' | 'verdant'

export interface SceneSummary {
  id: string
  title: string
  description: string
  act: string
  status: string
  actors: string[]
  extra: number
  tone?: SceneTone
  image?: string
  imageAlign?: string
  tags?: TagRef[]
}

export interface StatTile {
  key: string
  icon: string
  value: string
  delta?: string
  deltaTone?: 'up' | 'flat'
  accent: AccentToken
}
