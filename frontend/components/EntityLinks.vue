<script setup lang="ts">
const props = defineProps<{ subjectType: string; subjectId: string; subjectName?: string }>()

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()

interface Candidate {
  type: string
  id: string
  name: string
  image?: string
}
interface LinkRow {
  id: string
  kind: string
  other: Candidate
}

const links = ref<LinkRow[]>([])
const candidates = ref<Candidate[]>([])
const search = ref('')
const dragOver = ref(false)

const typeIcon: Record<string, string> = {
  character: 'lucide:user',
  item: 'lucide:gem',
  scene: 'lucide:castle',
  image: 'lucide:image',
}
const typeRoute: Record<string, string> = {
  character: '/characters/',
  item: '/items/',
  scene: '/scenes/',
  image: '/images/',
}

async function loadLinks() {
  links.value = await api<LinkRow[]>(`/campaigns/${campaign.currentId}/links?subjectId=${props.subjectId}`)
}

async function loadCandidates() {
  const cid = campaign.currentId
  const [chars, items, scenes, images] = await Promise.all([
    api<any[]>(`/campaigns/${cid}/characters`),
    api<any[]>(`/campaigns/${cid}/items`),
    api<any[]>(`/campaigns/${cid}/scenes`),
    api<any[]>(`/campaigns/${cid}/images`),
  ])
  candidates.value = [
    ...chars.map((x) => ({ type: 'character', id: x.id, name: x.name, image: x.imageUrl })),
    ...items.map((x) => ({ type: 'item', id: x.id, name: x.name, image: x.imageUrl })),
    ...scenes.map((x) => ({ type: 'scene', id: x.id, name: x.name, image: x.imageUrl })),
    ...images.map((x) => ({ type: 'image', id: x.id, name: x.name, image: x.imageUrl })),
  ]
}

onMounted(async () => {
  await campaign.ensure()
  await Promise.all([loadLinks(), loadCandidates()])
})

const linkedIds = computed(() => new Set(links.value.map((l) => l.other.id)))
const picker = computed(() =>
  candidates.value
    .filter((c) => !(c.type === props.subjectType && c.id === props.subjectId))
    .filter((c) => !linkedIds.value.has(c.id))
    .filter((c) => c.name.toLowerCase().includes(search.value.toLowerCase().trim()))
    .slice(0, 40),
)

async function addLink(cand: Candidate) {
  await api(`/campaigns/${campaign.currentId}/links`, {
    method: 'POST',
    body: { fromType: props.subjectType, fromId: props.subjectId, toType: cand.type, toId: cand.id },
  })
  await loadLinks()
}
async function removeLink(linkId: string) {
  await api(`/campaigns/${campaign.currentId}/links/${linkId}`, { method: 'DELETE' })
  await loadLinks()
}

function onDragStart(e: DragEvent, cand: Candidate) {
  e.dataTransfer?.setData('text/aw-entity', JSON.stringify(cand))
  if (e.dataTransfer) e.dataTransfer.effectAllowed = 'copy'
}
function onDrop(e: DragEvent) {
  dragOver.value = false
  const raw = e.dataTransfer?.getData('text/aw-entity')
  if (!raw) return
  try {
    const cand = JSON.parse(raw) as Candidate
    if (cand.type === props.subjectType && cand.id === props.subjectId) return
    if (!linkedIds.value.has(cand.id)) addLink(cand)
  } catch {
    /* ignore */
  }
}
</script>

<template>
  <section class="links">
    <div class="lh">
      <Icon name="lucide:link" />
      <h3>{{ t('links.title') }}</h3>
    </div>

    <div
      class="dropzone"
      :class="{ over: dragOver, empty: !links.length }"
      @dragover.prevent="dragOver = true"
      @dragleave="dragOver = false"
      @drop.prevent="onDrop"
    >
      <template v-if="links.length">
        <div v-for="l in links" :key="l.id" class="chip linked">
          <NuxtLink :to="`${typeRoute[l.other.type]}${l.other.id}`" class="cl">
            <span class="ci"><img v-if="l.other.image" :src="l.other.image" :alt="l.other.name" /><Icon v-else :name="typeIcon[l.other.type]" /></span>
            <span class="cn">{{ l.other.name }}</span>
          </NuxtLink>
          <button class="x" :title="t('links.remove')" @click="removeLink(l.id)"><Icon name="lucide:x" /></button>
        </div>
      </template>
      <span v-else class="hint"><Icon name="lucide:mouse-pointer-2" /> {{ t('links.dropHint') }}</span>
    </div>

    <div class="picker">
      <div class="psearch">
        <Icon name="lucide:search" />
        <input v-model="search" :placeholder="t('links.search')" />
      </div>
      <div v-if="picker.length" class="pchips">
        <button
          v-for="c in picker"
          :key="`${c.type}-${c.id}`"
          class="chip drag"
          draggable="true"
          @dragstart="onDragStart($event, c)"
          @click="addLink(c)"
        >
          <span class="ci sm"><Icon :name="typeIcon[c.type]" /></span>
          <span class="cn">{{ c.name }}</span>
          <Icon name="lucide:plus" class="add" />
        </button>
      </div>
      <p v-else class="empty">{{ t('links.empty') }}</p>
    </div>
  </section>
</template>

<style lang="scss" scoped>
.links { display: flex; flex-direction: column; gap: 12px; }
.lh {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 6px;
  :deep(svg) { width: 17px; height: 17px; color: var(--secondary); }
  h3 { font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; }
}

.dropzone {
  display: flex;
  flex-wrap: wrap;
  gap: 9px;
  min-height: 62px;
  padding: 14px;
  border-radius: 16px;
  border: 1.5px dashed var(--line-strong);
  background: var(--surface);
  transition: 0.2s;
  align-items: center;

  &.empty { justify-content: center; }
  &.over { border-color: var(--secondary); background: rgba(183, 104, 255, 0.08); box-shadow: var(--glow-secondary) inset; }
  .hint {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-family: var(--font-mono);
    font-size: 0.72rem;
    color: var(--ink-faint);
    :deep(svg) { width: 15px; height: 15px; }
  }
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border-radius: 10px;
  border: 1px solid var(--line-strong);
  background: var(--surface-2);
  font-size: 0.84rem;
  color: var(--ink);
  transition: 0.18s;

  .ci {
    width: 26px;
    height: 26px;
    border-radius: 7px;
    flex: none;
    display: grid;
    place-items: center;
    overflow: hidden;
    background: color-mix(in srgb, var(--secondary) 14%, transparent);
    color: var(--secondary);
    :deep(svg) { width: 14px; height: 14px; }
    img { width: 100%; height: 100%; object-fit: cover; }
    &.sm { width: 22px; height: 22px; :deep(svg) { width: 12px; height: 12px; } }
  }
  .cn { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 180px; }
}

.chip.linked {
  padding: 5px 6px 5px 6px;
  .cl { display: inline-flex; align-items: center; gap: 8px; text-decoration: none; color: inherit; }
  &:hover { border-color: var(--secondary); }
  .x {
    width: 22px;
    height: 22px;
    display: grid;
    place-items: center;
    border: 0;
    background: transparent;
    color: var(--ink-faint);
    border-radius: 6px;
    cursor: pointer;
    :deep(svg) { width: 13px; height: 13px; }
    &:hover { color: var(--ember); background: rgba(255, 106, 85, 0.12); }
  }
}

.picker { display: flex; flex-direction: column; gap: 10px; }
.psearch {
  display: flex;
  align-items: center;
  gap: 9px;
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 11px;
  padding: 9px 13px;
  :deep(svg) { width: 16px; height: 16px; color: var(--ink-faint); flex: none; }
  input { flex: 1; background: transparent; border: 0; outline: 0; color: var(--ink); font-family: var(--font-body); font-size: 0.88rem; }
}
.pchips { display: flex; flex-wrap: wrap; gap: 8px; max-height: 200px; overflow-y: auto; }
.chip.drag {
  padding: 6px 11px 6px 7px;
  cursor: grab;
  color: var(--ink-dim);
  &:hover { color: var(--ink); border-color: var(--secondary); box-shadow: var(--glow-secondary); }
  &:active { cursor: grabbing; }
  .add { width: 13px; height: 13px; color: var(--secondary); margin-left: 2px; }
}
.empty { font-size: 0.8rem; color: var(--ink-faint); font-family: var(--font-mono); }
</style>
