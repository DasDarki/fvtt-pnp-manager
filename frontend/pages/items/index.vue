<script setup lang="ts">
import type { Rarity, TagRef } from '~/types/entities'
import type { ApiItem, ApiEntityTag, ApiTag } from '~/types/api'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()
const ctl = useListControls('items')

useSeoMeta({ title: () => t('nav.items') })

const { data, pending } = useAsyncData(
  'items',
  async () => {
    await campaign.ensure()
    const cid = campaign.currentId
    const [rows, ets, tags] = await Promise.all([
      api<ApiItem[]>(`/campaigns/${cid}/items`),
      api<ApiEntityTag[]>(`/campaigns/${cid}/entity-tags?entityType=item`),
      api<ApiTag[]>(`/campaigns/${cid}/tags`),
    ])
    const tagMap: Record<string, TagRef[]> = {}
    for (const e of ets) (tagMap[e.entityId] ??= []).push({ id: e.tag.id, name: e.tag.name, color: e.tag.color })
    return { items: rows.map((r) => ({ ...toItemVM(r), tags: tagMap[r.id] || [] })), tags }
  },
  { server: false, default: () => ({ items: [], tags: [] as ApiTag[] }) },
)

const rarities: Rarity[] = ['common', 'uncommon', 'rare', 'epic', 'legendary', 'artifact']
const rarityOptions = computed(() => [
  { value: 'all', label: t('filter.all') },
  ...rarities.map((r) => ({ value: r, label: t(`rarity.${r}`) })),
])
const sortOptions = computed(() => [
  { value: 'recent', label: t('sort.recent') },
  { value: 'name', label: t('sort.name') },
])

const filtered = computed(() => {
  const q = ctl.search.toLowerCase().trim()
  let list = data.value.items.filter(
    (i) =>
      (ctl.primary === 'all' || i.rarity === ctl.primary) &&
      i.name.toLowerCase().includes(q) &&
      (!ctl.hasImage || !!i.image) &&
      (ctl.tags.length === 0 || ctl.tags.every((tid) => (i.tags || []).some((tg) => tg.id === tid))),
  )
  if (ctl.sort === 'name') list = [...list].sort((a, b) => a.name.localeCompare(b.name))
  return list
})

async function create() {
  if (!campaign.currentId) return
  const it = await api<ApiItem>(`/campaigns/${campaign.currentId}/items`, {
    method: 'POST',
    body: { name: 'Neues Item', itemType: 'Gegenstand', rarity: 'common', summary: 'Frisch erstellt', systemData: { icon: 'lucide:gem' } },
  })
  await navigateTo(`/items/${it.id}`)
}
</script>

<template>
  <div>
    <PageHeader
      v-model:search="ctl.search"
      :title="t('nav.items')"
      icon="lucide:gem"
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
      v-model:primary="ctl.primary"
      v-model:tags="ctl.tags"
      v-model:sort="ctl.sort"
      v-model:hasImage="ctl.hasImage"
      :primary-options="rarityOptions"
      :sort-options="sortOptions"
      :campaign-tags="data.tags"
      show-image-filter
    />

    <EntityGrid :empty="!pending && filtered.length === 0" :min="240">
      <NuxtLink v-for="i in filtered" :key="i.id" :to="`/items/${i.id}`" class="card-link">
        <ItemCard :item="i" />
      </NuxtLink>
      <template #empty>{{ t('empty.items') }}</template>
    </EntityGrid>
  </div>
</template>

<style lang="scss" scoped>
.card-link { display: block; text-decoration: none; color: inherit; }
@media (max-width: 820px) {
  .hide-m { display: none; }
}
</style>
