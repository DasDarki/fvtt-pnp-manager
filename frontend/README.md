# Aetherwright — Frontend

Nuxt 3 · Tailwind · SCSS · Pinia · vue-i18n. Teil der [-wright-Suite](../umsetzungsplan.md).

## Setup (bun)

```bash
bun install
bun run dev      # http://localhost:3000
```

Build:

```bash
bun run build
bun run preview
```

## Stack & Konventionen

- **Design-Tokens** als CSS-Variablen in `assets/scss/main.scss` (Dark default, Light über
  `html[data-theme="light"]`), in Tailwind gespiegelt (`tailwind.config.ts`).
- **Fonts** werden via `@nuxt/fonts` zur Build-Zeit heruntergeladen und lokal gebundlet
  (Cinzel, Cinzel Decorative, Manrope, JetBrains Mono) — kein Google-CDN zur Laufzeit.
- **i18n** (`@nuxtjs/i18n`): DE Default, EN optional, `strategy: no_prefix`. Strings unter
  `i18n/locales/*.json`.
- **SEO/Head**: zentral in `app.vue` über `useHead` (Title-Template) + `useSeoMeta`
  (OG/Twitter, locale-abhängig). Marken-Konstanten in `app.config.ts`.
- **Theme**: `composables/useTheme.ts` (persistiert in `localStorage`,
  No-Flash-Inline-Script in `nuxt.config.ts`).

> OG-Bild liegt als SVG (`public/og/aetherwright-og.svg`). Für maximale Plattform-
> Kompatibilität später als PNG rastern oder `nuxt-og-image` ergänzen.
