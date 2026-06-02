export const TAG_COLORS = ['primary', 'secondary', 'magenta', 'gold', 'emerald', 'ember']

export function tagColorFor(name: string): string {
  let h = 0
  for (const ch of name) h = (h * 31 + ch.charCodeAt(0)) >>> 0
  return TAG_COLORS[h % TAG_COLORS.length]
}
