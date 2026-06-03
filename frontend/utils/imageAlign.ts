export type ImageAlign = 'top' | 'center' | 'bottom'

export const IMAGE_ALIGNS: ImageAlign[] = ['top', 'center', 'bottom']

export function objPos(align?: string): string {
  if (align === 'top') return 'center top'
  if (align === 'bottom') return 'center bottom'
  return 'center center'
}

export function alignIcon(align: ImageAlign): string {
  if (align === 'top') return 'lucide:align-vertical-justify-start'
  if (align === 'bottom') return 'lucide:align-vertical-justify-end'
  return 'lucide:align-vertical-justify-center'
}
