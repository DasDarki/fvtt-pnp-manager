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
