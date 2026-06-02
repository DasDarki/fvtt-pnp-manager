<script setup lang="ts">
import type { TagRef } from '~/types/entities'
import type { ApiImage, ApiEntityTag, ApiTag } from '~/types/api'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()
const ctl = useListControls('images')

useSeoMeta({ title: () => t('nav.images') })

const { data, pending } = useAsyncData(
  'images',
  async () => {
    await campaign.ensure()
    const cid = campaign.currentId
    const [rows, ets, tags] = await Promise.all([
      api<ApiImage[]>(`/campaigns/${cid}/images`),
      api<ApiEntityTag[]>(`/campaigns/${cid}/entity-tags?entityType=image`),
      api<ApiTag[]>(`/campaigns/${cid}/tags`),
    ])
    const tagMap: Record<string, TagRef[]> = {}
    for (const e of ets) (tagMap[e.entityId] ??= []).push({ id: e.tag.id, name: e.tag.name, color: e.tag.color })
    return { images: rows.map((r) => ({ ...toImageVM(r), tags: tagMap[r.id] || [] })), tags }
  },
  { server: false, default: () => ({ images: [], tags: [] as ApiTag[] }) },
)

const sortOptions = computed(() => [
  { value: 'recent', label: t('sort.recent') },
  { value: 'name', label: t('sort.name') },
])

const filtered = computed(() => {
  const q = ctl.search.toLowerCase().trim()
  let list = data.value.images.filter(
    (i) =>
      i.name.toLowerCase().includes(q) &&
      (!ctl.hasImage || !!i.image) &&
      (ctl.tags.length === 0 || ctl.tags.every((tid) => (i.tags || []).some((tg) => tg.id === tid))),
  )
  if (ctl.sort === 'name') list = [...list].sort((a, b) => a.name.localeCompare(b.name))
  return list
})

async function create() {
  if (!campaign.currentId) return
  const created = await api<ApiImage>(`/campaigns/${campaign.currentId}/images`, {
    method: 'POST',
    body: { name: 'Neues Bild', pushAs: 'empty_actor', notes: '' },
  })
  await navigateTo(`/images/${created.id}`)
}
</script>

<template>
  <div>
    <PageHeader
      v-model:search="ctl.search"
      :title="t('nav.images')"
      icon="lucide:image"
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
      :sort-options="sortOptions"
      :campaign-tags="data.tags"
    />

    <EntityGrid :empty="!pending && filtered.length === 0" :min="240">
      <NuxtLink v-for="i in filtered" :key="i.id" :to="`/images/${i.id}`" class="card-link">
        <ImageCard :image="i" />
      </NuxtLink>
      <template #empty>{{ t('image.empty') }}</template>
    </EntityGrid>
  </div>
</template>

<style lang="scss" scoped>
.card-link { display: block; text-decoration: none; color: inherit; }
@media (max-width: 820px) {
  .hide-m { display: none; }
}
</style>
