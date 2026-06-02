export interface StylePreset {
  id: string
  label: string
  prompt: string
}

export const STYLE_PRESETS: StylePreset[] = [
  { id: 'dark-fantasy', label: 'Dark Fantasy', prompt: 'dark fantasy, dramatic chiaroscuro lighting, painterly, intricate detail, muted palette' },
  { id: 'gothic-horror', label: 'Gothic Horror', prompt: 'gothic horror, candlelit, desaturated, ominous atmosphere, oil painting' },
  { id: 'high-fantasy', label: 'High Fantasy', prompt: 'high fantasy illustration, vibrant, epic, golden hour, richly detailed' },
  { id: 'oil-painting', label: 'Ölgemälde', prompt: 'classical oil painting, visible brushstrokes, warm varnish, renaissance lighting' },
  { id: 'ink-noir', label: 'Tinte & Noir', prompt: 'high-contrast ink illustration, noir, bold blacks, dramatic shadow' },
  { id: 'watercolor', label: 'Aquarell', prompt: 'soft watercolor, loose washes, delicate edges, paper texture' },
  { id: 'anime', label: 'Anime', prompt: 'anime style, cel shaded, clean linework, expressive, cinematic composition' },
  { id: 'grim-realism', label: 'Grimmiger Realismus', prompt: 'gritty cinematic realism, volumetric light, film grain, photoreal detail' },
]

export function presetById(id?: string | null): StylePreset | undefined {
  return STYLE_PRESETS.find((p) => p.id === id)
}
