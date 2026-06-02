const PUBLIC = new Set(['/', '/login'])

export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) return

  const auth = useAuthStore()

  if (!auth.isAuthed && !PUBLIC.has(to.path)) {
    return navigateTo('/login')
  }
  if (auth.isAuthed && to.path === '/login') {
    return navigateTo('/dashboard')
  }
})
