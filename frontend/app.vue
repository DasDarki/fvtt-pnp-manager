<script setup lang="ts">
const { locale, t } = useI18n()
const appConfig = useAppConfig()
const requestUrl = useRequestURL()

const fullTitle = computed(() => `${appConfig.site.name} — ${t('meta.tagline')}`)
const ogLocale = computed(() => (locale.value === 'de' ? 'de_DE' : 'en_US'))
const ogImageUrl = computed(() => appConfig.site.url + appConfig.site.ogImage)

useHead({
  titleTemplate: (titleChunk?: string) =>
    titleChunk ? `${titleChunk} · ${appConfig.site.name}` : fullTitle.value,
  link: [{ rel: 'canonical', href: () => requestUrl.href }],
})

useSeoMeta({
  description: () => t('meta.description'),
  ogType: 'website',
  ogSiteName: appConfig.site.name,
  ogTitle: () => fullTitle.value,
  ogDescription: () => t('meta.description'),
  ogImage: () => ogImageUrl.value,
  ogImageAlt: () => fullTitle.value,
  ogUrl: () => requestUrl.href,
  ogLocale: () => ogLocale.value,
  twitterCard: 'summary_large_image',
  twitterSite: appConfig.site.twitter,
  twitterTitle: () => fullTitle.value,
  twitterDescription: () => t('meta.description'),
  twitterImage: () => ogImageUrl.value,
})
</script>

<template>
  <div>
    <NuxtRouteAnnouncer />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>
