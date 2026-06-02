export function useListControls(key: string) {
  const state = reactive({
    search: '',
    sort: 'recent',
    primary: 'all',
    tags: [] as string[],
    hasImage: false,
  })

  if (import.meta.client) {
    const saved = localStorage.getItem('aw-list-' + key)
    if (saved) {
      try {
        const d = JSON.parse(saved)
        if (d.sort) state.sort = d.sort
        if (d.primary) state.primary = d.primary
        if (typeof d.hasImage === 'boolean') state.hasImage = d.hasImage
      } catch {
        /* ignore */
      }
    }
    watch(
      () => [state.sort, state.primary, state.hasImage],
      () => {
        localStorage.setItem(
          'aw-list-' + key,
          JSON.stringify({ sort: state.sort, primary: state.primary, hasImage: state.hasImage }),
        )
      },
    )
  }

  return state
}
