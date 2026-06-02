export interface ProviderMeta {
  label: string
  defaultModel: string
  icon: string
  keyHost: string
}

export const PROVIDER_META: Record<string, ProviderMeta> = {
  openai: { label: 'OpenAI · DALL·E / GPT-Image', defaultModel: 'dall-e-3', icon: 'lucide:sparkles', keyHost: 'platform.openai.com' },
  stability: { label: 'Stability AI', defaultModel: 'core', icon: 'lucide:layers', keyHost: 'platform.stability.ai' },
  replicate: { label: 'Replicate · Flux', defaultModel: 'black-forest-labs/flux-schnell', icon: 'lucide:boxes', keyHost: 'replicate.com' },
  gemini: { label: 'Google · Imagen', defaultModel: 'imagen-3.0-generate-002', icon: 'lucide:gem', keyHost: 'aistudio.google.com' },
}

export function providerLabel(provider: string): string {
  return PROVIDER_META[provider]?.label || provider
}

export function providerIcon(provider: string): string {
  return PROVIDER_META[provider]?.icon || 'lucide:cpu'
}
