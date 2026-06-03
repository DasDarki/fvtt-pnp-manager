export default defineNuxtConfig({
  compatibilityDate: '2025-01-15',
  devtools: { enabled: true },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@nuxtjs/i18n',
    '@nuxt/fonts',
    '@nuxt/icon',
  ],

  css: ['~/assets/scss/main.scss'],

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1',
    },
  },

  icon: {
    mode: 'svg',
    serverBundle: 'local',
    // Eigener Server-Endpoint NICHT unter /api (das routet der Proxy aufs Backend).
    localApiEndpoint: '/_nuxt_icon',
    // Alle genutzten Icons in den Client-Bundle → kein Runtime-Fetch, kein Hydration-Mismatch.
    clientBundle: {
      scan: true,
      sizeLimitKb: 512,
      icons: [
        'lucide:activity', 'lucide:arrow-down-up', 'lucide:arrow-left', 'lucide:arrow-right-left',
        'lucide:boxes', 'lucide:castle', 'lucide:check', 'lucide:check-check', 'lucide:chevrons-up-down',
        'lucide:circle-alert', 'lucide:copy', 'lucide:cpu', 'lucide:crown', 'lucide:external-link',
        'lucide:flame', 'lucide:flask-conical', 'lucide:folder', 'lucide:folder-sync', 'lucide:folder-tree',
        'lucide:gem', 'lucide:ghost', 'lucide:hexagon', 'lucide:image', 'lucide:image-up', 'lucide:images',
        'lucide:folder-plus', 'lucide:upload', 'lucide:key-round', 'lucide:layers',
        'lucide:layout-dashboard', 'lucide:library', 'lucide:link', 'lucide:list', 'lucide:loader-circle', 'lucide:log-in',
        'lucide:key-square', 'lucide:puzzle',
        'lucide:log-out', 'lucide:menu', 'lucide:moon-star', 'lucide:mouse-pointer-2', 'lucide:palette',
        'lucide:pin', 'lucide:plug-zap', 'lucide:plus', 'lucide:refresh-cw', 'lucide:save', 'lucide:scroll',
        'lucide:scroll-text', 'lucide:search', 'lucide:search-x', 'lucide:settings-2', 'lucide:sparkles',
        'lucide:star', 'lucide:sun', 'lucide:sword', 'lucide:tags', 'lucide:trash', 'lucide:trash-2',
        'lucide:trending-up', 'lucide:triangle-alert', 'lucide:user', 'lucide:user-plus', 'lucide:users',
        'lucide:wand-2', 'lucide:wand-sparkles', 'lucide:x',
      ],
    },
  },

  typescript: {
    strict: true,
  },

  fonts: {
    families: [
      { name: 'Cinzel', provider: 'google', weights: [500, 600, 700] },
      { name: 'Cinzel Decorative', provider: 'google', weights: [700, 900] },
      { name: 'Manrope', provider: 'google', weights: [300, 400, 500, 600, 700, 800] },
      { name: 'JetBrains Mono', provider: 'google', weights: [400, 500, 600] },
    ],
  },

  i18n: {
    defaultLocale: 'de',
    strategy: 'no_prefix',
    langDir: 'locales',
    locales: [
      { code: 'de', language: 'de-DE', name: 'Deutsch', file: 'de.json' },
      { code: 'en', language: 'en-US', name: 'English', file: 'en.json' },
    ],
    detectBrowserLanguage: false,
    bundle: { optimizeTranslationDirective: false },
  },

  app: {
    head: {
      htmlAttrs: { lang: 'de', 'data-theme': 'dark' },
      charset: 'utf-8',
      viewport: 'width=device-width, initial-scale=1, viewport-fit=cover',
      link: [{ rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }],
      meta: [
        { name: 'theme-color', content: '#07060d' },
        { name: 'color-scheme', content: 'dark light' },
      ],
      script: [
        {
          innerHTML:
            "try{var t=localStorage.getItem('aetherwright-theme');if(t){document.documentElement.dataset.theme=t}}catch(e){}",
          tagPosition: 'head',
        },
      ],
    },
  },
})
