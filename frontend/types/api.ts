export interface ApiUser {
  id: string
  email: string
  name: string
}

export interface AuthSession {
  accessToken: string
  refreshToken: string
  user: ApiUser
}

export interface ApiCampaign {
  id: string
  name: string
  slug: string
  description: string
  ruleset: string
  forgeRootPath: string
  artStyle: string
  stylePrompt: string
}

export interface DalleResult {
  id: string
  prompt: string
  status: string
  size: string
  imageUrl: string
  mock: boolean
  error: string
  createdAt: string
}

export interface ApiAsset {
  id: string
  url: string
  name: string
  mime: string
  kind: string
  source: string
  prompt: string
  createdAt: string
}

export interface ApiProvider {
  provider: string
  model: string
  hasKey: boolean
}

export interface ProvidersResponse {
  providers: ApiProvider[]
  catalog: string[]
}

export interface ApiCharacter {
  id: string
  campaignId: string
  name: string
  summary: string
  characterType: string
  status: string
  systemData: Record<string, any> | null
  imageUrl?: string
  syncState?: string
  foundryUuid?: string
  folderId?: string | null
}

export interface ApiItem {
  id: string
  name: string
  summary: string
  itemType: string
  rarity: string
  attuned: boolean
  systemData: Record<string, any> | null
  imageUrl?: string
  syncState?: string
  foundryUuid?: string
  folderId?: string | null
}

export interface ApiImage {
  id: string
  name: string
  notes: string
  pushAs: string
  assetId?: string | null
  imageUrl?: string
  syncState?: string
  foundryUuid?: string
  folderId?: string | null
}

export interface ApiScene {
  id: string
  name: string
  summary: string
  sceneStatus: string
  systemData: Record<string, any> | null
  imageUrl?: string
  syncState?: string
  foundryUuid?: string
  folderId?: string | null
}

export interface ApiMemory {
  id: string
  title: string
  body: string
  kind: string
  level: string
  subjectType: string
  subjectId: string | null
  subjectLabel: string
  acknowledged: boolean
  pinned: boolean
  createdAt: string
}

export interface ApiTag {
  id: string
  name: string
  color: string
}

export interface ApiEntityTag {
  id: string
  entityId: string
  entityType: string
  tag: ApiTag
}
