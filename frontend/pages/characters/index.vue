<script setup lang="ts">
import type { CharacterStatus, TagRef } from '~/types/entities'
import type { ApiCharacter, ApiEntityTag, ApiTag } from '~/types/api'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()
const ctl = useListControls('characters')

useSeoMeta({ title: () => t('nav.characters') })

const { data, pending, refresh } = useAsyncData(
  'characters',
  async () => {
    await campaign.ensure()
    const cid = campaign.currentId
    const [rows, ets, tags] = await Promise.all([
      api<ApiCharacter[]>(`/campaigns/${cid}/characters`),
      api<ApiEntityTag[]>(`/campaigns/${cid}/entity-tags?entityType=character`),
      api<ApiTag[]>(`/campaigns/${cid}/tags`),
    ])
    const tagMap: Record<string, TagRef[]> = {}
    for (const e of ets) (tagMap[e.entityId] ??= []).push({ id: e.tag.id, name: e.tag.name, color: e.tag.color })
    return {
      characters: rows.map((r) => ({ ...toCharacterVM(r), tags: tagMap[r.id] || [] })),
      tags,
    }
  },
  { server: false, default: () => ({ characters: [], tags: [] as ApiTag[] }) },
)

const statusOptions = computed(() => [
  { value: 'all', label: t('filter.all') },
  ...(['alive', 'dead', 'hunted', 'unknown'] as CharacterStatus[]).map((s) => ({ value: s, label: t(`status.${s}`) })),
])
const sortOptions = computed(() => [
  { value: 'recent', label: t('sort.recent') },
  { value: 'name', label: t('sort.name') },
])

const filtered = computed(() => {
  const q = ctl.search.toLowerCase().trim()
  let list = data.value.characters.filter(
    (c) =>
      (ctl.primary === 'all' || c.status === ctl.primary) &&
      c.name.toLowerCase().includes(q) &&
      (!ctl.hasImage || !!c.image) &&
      (ctl.tags.length === 0 || ctl.tags.every((tid) => (c.tags || []).some((tg) => tg.id === tid))),
  )
  if (ctl.sort === 'name') list = [...list].sort((a, b) => a.name.localeCompare(b.name))
  return list
})

const creating = ref(false)
async function create() {
  if (!campaign.currentId) return
  creating.value = true
  try {
    const ch = await api<ApiCharacter>(`/campaigns/${campaign.currentId}/characters`, {
      method: 'POST',
      body: {
        name: 'Neuer Charakter',
        characterType: 'npc',
        status: 'unknown',
        systemData: {
          subtitle: 'Frisch erstellt · NSC',
          proficiency: 2,
          abilities: { str: 10, dex: 10, con: 10, int: 10, wis: 10, cha: 10 },
          ac: 10,
          hp: 8,
          hpMax: 8,
          level: 1,
        },
      },
    })
    await navigateTo(`/characters/${ch.id}`)
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <div>
    <PageHeader
      v-model:search="ctl.search"
      :title="t('nav.characters')"
      icon="lucide:users"
      :count="filtered.length"
      :count-label="t('common.entries')"
      :search-placeholder="t('search.entity')"
    >
      <template #actions>
        <AwButton icon="lucide:user-plus" @click="create">
          <span class="hide-m">{{ t('actions.create') }}</span>
        </AwButton>
      </template>
    </PageHeader>

    <FilterBar
      v-model:primary="ctl.primary"
      v-model:tags="ctl.tags"
      v-model:sort="ctl.sort"
      v-model:hasImage="ctl.hasImage"
      :primary-options="statusOptions"
      :sort-options="sortOptions"
      :campaign-tags="data.tags"
      show-image-filter
    />

    <EntityGrid :empty="!pending && filtered.length === 0">
      <NuxtLink v-for="c in filtered" :key="c.id" :to="`/characters/${c.id}`" class="card-link">
        <CharacterCard :character="c" />
      </NuxtLink>
      <template #empty>{{ t('empty.characters') }}</template>
    </EntityGrid>
  </div>
</template>

<style lang="scss" scoped>
.card-link { display: block; text-decoration: none; color: inherit; }
@media (max-width: 820px) {
  .hide-m { display: none; }
}
</style>
