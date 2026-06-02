export type Theme = 'dark' | 'light'

const STORAGE_KEY = 'aetherwright-theme'

export function useTheme() {
  const theme = useState<Theme>('aw-theme', () => 'dark')

  function apply(next: Theme) {
    theme.value = next
    if (import.meta.client) {
      localStorage.setItem(STORAGE_KEY, next)
      document.documentElement.dataset.theme = next
    }
  }

  function toggle() {
    apply(theme.value === 'dark' ? 'light' : 'dark')
  }

  onMounted(() => {
    const stored = localStorage.getItem(STORAGE_KEY) as Theme | null
    if (stored === 'dark' || stored === 'light') {
      theme.value = stored
      document.documentElement.dataset.theme = stored
    }
  })

  return { theme, toggle, apply }
}
