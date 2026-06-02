<script setup lang="ts">
import type { ApiEntityTag, ApiTag } from '~/types/api'

const props = defineProps<{ subjectType: string; subjectId: string }>()

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()

const entityTags = ref<ApiEntityTag[]>([])
const allTags = ref<ApiTag[]>([])
const input = ref('')
const focused = ref(false)

async function load() {
  const cid = campaign.currentId
  const [et, tags] = await Promise.all([
    api<ApiEntityTag[]>(`/campaigns/${cid}/entity-tags?subjectId=${props.subjectId}`),
    api<ApiTag[]>(`/campaigns/${cid}/tags`),
  ])
  entityTags.value = et
  allTags.value = tags
}

onMounted(async () => {
  await campaign.ensure()
  await load()
})

const attachedIds = computed(() => new Set(entityTags.value.map((e) => e.tag.id)))
const suggestions = computed(() => {
  const q = input.value.toLowerCase().trim()
  return allTags.value
    .filter((tg) => !attachedIds.value.has(tg.id))
    .filter((tg) => !q || tg.name.toLowerCase().includes(q))
    .slice(0, 6)
})

async function attach(tagId: string) {
  await api(`/campaigns/${campaign.currentId}/entity-tags`, {
    method: 'POST',
    body: { tagId, entityType: props.subjectType, entityId: props.subjectId },
  })
  input.value = ''
  await load()
}

async function createAndAttach() {
  const name = input.value.trim()
  if (!name) return
  const existing = allTags.value.find((tg) => tg.name.toLowerCase() === name.toLowerCase())
  if (existing) return attach(existing.id)
  const tag = await api<ApiTag>(`/campaigns/${campaign.currentId}/tags`, {
    method: 'POST',
    body: { name, color: tagColorFor(name) },
  })
  await attach(tag.id)
}

async function detach(etId: string) {
  await api(`/campaigns/${campaign.currentId}/entity-tags/${etId}`, { method: 'DELETE' })
  await load()
}
</script>

<template>
  <section class="tags">
    <div class="th">
      <Icon name="lucide:tags" />
      <h3>{{ t('tag.title') }}</h3>
    </div>

    <div class="row">
      <span v-for="e in entityTags" :key="e.id" class="chip" :style="{ '--c': `var(--${e.tag.color})` }">
        {{ e.tag.name }}
        <button class="x" :title="t('tag.remove')" @click="detach(e.id)"><Icon name="lucide:x" /></button>
      </span>

      <div class="add">
        <input
          v-model="input"
          class="tinp"
          :placeholder="t('tag.add')"
          @focus="focused = true"
          @blur="focused = false"
          @keydown.enter.prevent="createAndAttach"
        />
        <div v-if="focused && input.trim()" class="sug">
          <button v-for="s in suggestions" :key="s.id" class="srow" @mousedown.prevent="attach(s.id)">
            <span class="d" :style="{ background: `var(--${s.color})` }" /> {{ s.name }}
          </button>
          <button class="srow new" @mousedown.prevent="createAndAttach">
            <Icon name="lucide:plus" /> {{ t('tag.createNew', { name: input.trim() }) }}
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<style lang="scss" scoped>
.tags { display: flex; flex-direction: column; gap: 10px; }
.th {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 6px;
  :deep(svg) { width: 17px; height: 17px; color: var(--gold); }
  h3 { font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; }
}

.row { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; }

.chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 0.7rem;
  color: var(--c, var(--secondary));
  border: 1px solid color-mix(in srgb, var(--c, var(--secondary)) 45%, transparent);
  background: color-mix(in srgb, var(--c, var(--secondary)) 12%, transparent);
  padding: 5px 5px 5px 11px;
  border-radius: 999px;

  .x {
    width: 18px;
    height: 18px;
    display: grid;
    place-items: center;
    border: 0;
    background: transparent;
    color: inherit;
    opacity: 0.6;
    border-radius: 50%;
    cursor: pointer;
    :deep(svg) { width: 12px; height: 12px; }
    &:hover { opacity: 1; background: color-mix(in srgb, var(--c, var(--secondary)) 25%, transparent); }
  }
}

.add { position: relative; }
.tinp {
  font-family: var(--font-body);
  font-size: 0.82rem;
  color: var(--ink);
  background: var(--surface-2);
  border: 1px dashed var(--line-strong);
  border-radius: 999px;
  padding: 7px 14px;
  width: 160px;
  transition: 0.2s;
  &:focus { outline: 0; border-style: solid; border-color: var(--gold); box-shadow: 0 0 16px -6px var(--gold); }
}

.sug {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  z-index: 20;
  min-width: 200px;
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 12px;
  box-shadow: 0 24px 50px -22px #000;
  padding: 5px;
  overflow: hidden;

  .srow {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 9px;
    padding: 8px 10px;
    border: 0;
    background: transparent;
    color: var(--ink-dim);
    font-size: 0.82rem;
    border-radius: 8px;
    cursor: pointer;
    text-align: left;
    &:hover { background: var(--surface); color: var(--ink); }
    .d { width: 8px; height: 8px; border-radius: 50%; flex: none; }
    :deep(svg) { width: 13px; height: 13px; color: var(--gold); }
    &.new { color: var(--gold); border-top: 1px solid var(--line); margin-top: 3px; }
  }
}
</style>
