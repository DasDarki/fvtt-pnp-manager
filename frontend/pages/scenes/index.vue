<script setup lang="ts">
import type { TagRef } from '~/types/entities'
import type { ApiScene, ApiEntityTag, ApiTag } from '~/types/api'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()
const ctl = useListControls('scenes')

useSeoMeta({ title: () => t('nav.scenes') })

const { data, pending } = useAsyncData(
  'scenes',
  async () => {
    await campaign.ensure()
    const cid = campaign.currentId
    const [rows, ets, tags] = await Promise.all([
      api<ApiScene[]>(`/campaigns/${cid}/scenes`),
      api<ApiEntityTag[]>(`/campaigns/${cid}/entity-tags?entityType=scene`),
      api<ApiTag[]>(`/campaigns/${cid}/tags`),
    ])
    const tagMap: Record<string, TagRef[]> = {}
    for (const e of ets) (tagMap[e.entityId] ??= []).push({ id: e.tag.id, name: e.tag.name, color: e.tag.color })
    return { scenes: rows.map((r) => ({ ...toSceneVM(r), tags: tagMap[r.id] || [] })), tags }
  },
  { server: false, default: () => ({ scenes: [], tags: [] as ApiTag[] }) },
)

const sortOptions = computed(() => [
  { value: 'recent', label: t('sort.recent') },
  { value: 'name', label: t('sort.name') },
])

const filtered = computed(() => {
  const q = ctl.search.toLowerCase().trim()
  let list = data.value.scenes.filter(
    (s) =>
      s.title.toLowerCase().includes(q) &&
      (!ctl.hasImage || !!s.image) &&
      (ctl.tags.length === 0 || ctl.tags.every((tid) => (s.tags || []).some((tg) => tg.id === tid))),
  )
  if (ctl.sort === 'name') list = [...list].sort((a, b) => a.title.localeCompare(b.title))
  return list
})

async function create() {
  if (!campaign.currentId) return
  const sc = await api<ApiScene>(`/campaigns/${campaign.currentId}/scenes`, {
    method: 'POST',
    body: { name: 'Neue Szene', summary: 'Frisch erstellt', sceneStatus: 'draft', systemData: { act: 'Entwurf', status: 'Entwurf', actors: [], extra: 0, tone: 'arcane' } },
  })
  await navigateTo(`/scenes/${sc.id}`)
}
</script>

<template>
  <div>
    <PageHeader
      v-model:search="ctl.search"
      :title="t('nav.scenes')"
      icon="lucide:castle"
      :count="filtered.length"
      :count-label="t('common.entries')"
      :search-placeholder="t('search.entity')"
    >
      <template #actions>
        <AwButton icon="lucide:plus" @click="create">
          <span class="hide-m">{{ t('actions.create') }}</span>
        </AwButton>
      </template>
    </PageHeader>

    <FilterBar
      v-model:tags="ctl.tags"
      v-model:sort="ctl.sort"
      v-model:hasImage="ctl.hasImage"
      :sort-options="sortOptions"
      :campaign-tags="data.tags"
      show-image-filter
    />

    <EntityGrid :empty="!pending && filtered.length === 0" :min="360">
      <NuxtLink v-for="s in filtered" :key="s.id" :to="`/scenes/${s.id}`" class="card-link">
        <SceneCard :scene="s" />
      </NuxtLink>
      <template #empty>{{ t('empty.scenes') }}</template>
    </EntityGrid>
  </div>
</template>

<style lang="scss" scoped>
.card-link { display: block; text-decoration: none; color: inherit; }
@media (max-width: 820px) {
  .hide-m { display: none; }
}
</style>
