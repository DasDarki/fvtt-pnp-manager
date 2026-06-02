import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', () => {
  const paletteOpen = ref(false)
  const sidebarOpen = ref(false)

  function openPalette() {
    paletteOpen.value = true
  }
  function closePalette() {
    paletteOpen.value = false
  }
  function togglePalette() {
    paletteOpen.value = !paletteOpen.value
  }
  function toggleSidebar() {
    sidebarOpen.value = !sidebarOpen.value
  }
  function closeSidebar() {
    sidebarOpen.value = false
  }

  return {
    paletteOpen,
    sidebarOpen,
    openPalette,
    closePalette,
    togglePalette,
    toggleSidebar,
    closeSidebar,
  }
})
