<script setup lang="ts">
import type { MemoryLevel } from '~/types/entities'
import type { ApiMemory } from '~/types/api'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()

useSeoMeta({ title: () => t('nav.memories') })

const { data: memories, pending, refresh } = useAsyncData(
  'memories-page',
  async () => {
    await campaign.ensure()
    const rows = await api<ApiMemory[]>(`/campaigns/${campaign.currentId}/memories`)
    return rows.map(toMemoryVM)
  },
  { server: false, default: () => [] },
)

const search = ref('')
const level = ref<'all' | MemoryLevel>('all')
const levels = ['all', 'info', 'notice', 'warning', 'critical'] as const
const showComposer = ref(false)

const filtered = computed(() =>
  (memories.value || []).filter(
    (m) =>
      (level.value === 'all' || m.level === level.value) &&
      (m.title.toLowerCase().includes(search.value.toLowerCase().trim()) ||
        m.body.toLowerCase().includes(search.value.toLowerCase().trim())),
  ),
)

async function ack(id: string) {
  await api(`/campaigns/${campaign.currentId}/memories/${id}`, { method: 'PATCH', body: { acknowledged: true } })
  await refresh()
}
async function togglePin(m: { id: string; pinned?: boolean }) {
  await api(`/campaigns/${campaign.currentId}/memories/${m.id}`, { method: 'PATCH', body: { pinned: !m.pinned } })
  await refresh()
}
async function del(id: string) {
  await api(`/campaigns/${campaign.currentId}/memories/${id}`, { method: 'DELETE' })
  await refresh()
}

async function onSaved() {
  showComposer.value = false
  await refresh()
}
</script>

<template>
  <div>
    <PageHeader
      v-model:search="search"
      :title="t('nav.memories')"
      icon="lucide:sparkles"
      :count="filtered.length"
      :count-label="t('common.entries')"
      :search-placeholder="t('search.entity')"
    >
      <template #actions>
        <AwButton icon="lucide:plus" @click="showComposer = !showComposer">
          <span class="hide-m">{{ t('memory.new') }}</span>
        </AwButton>
      </template>
    </PageHeader>

    <MemoryComposer v-if="showComposer" subject-type="campaign" :on-saved="onSaved" class="mb" />

    <div class="chips">
      <button
        v-for="l in levels"
        :key="l"
        class="chip"
        :class="{ on: level === l }"
        @click="level = l"
      >
        {{ l === 'all' ? t('filter.all') : t(`level.${l}`) }}
      </button>
    </div>

    <div v-if="filtered.length" class="list">
      <MemoryCard
        v-for="m in filtered"
        :key="m.id"
        :memory="m"
        :on-ack="() => ack(m.id)"
        :on-pin="() => togglePin(m)"
        :on-delete="() => del(m.id)"
      />
    </div>
    <div v-else-if="!pending" class="empty">
      <Icon name="lucide:sparkles" />
      {{ t('memory.empty') }}
    </div>
  </div>
</template>

<style lang="scss" scoped>
.mb { margin-bottom: 18px; }
.chips {
  display: flex;
  gap: 9px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}
.chip {
  font-family: var(--font-mono);
  font-size: 0.68rem;
  letter-spacing: 0.06em;
  color: var(--ink-dim);
  padding: 7px 14px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: var(--surface);
  cursor: pointer;
  transition: 0.2s;
  &:hover { color: var(--ink); border-color: var(--line-strong); }
  &.on { color: #06040c; background: var(--grad-arcane); border-color: transparent; box-shadow: var(--glow-primary); }
}
:global(html[data-theme='light']) .chip.on { color: #fff; }

.list { display: flex; flex-direction: column; gap: 12px; }

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding: 70px 20px;
  text-align: center;
  color: var(--ink-faint);
  font-family: var(--font-mono);
  font-size: 0.82rem;
  border: 1.5px dashed var(--line-strong);
  border-radius: 22px;
  :deep(svg) { width: 32px; height: 32px; opacity: 0.7; }
}

@media (max-width: 820px) {
  .hide-m { display: none; }
}
</style>
