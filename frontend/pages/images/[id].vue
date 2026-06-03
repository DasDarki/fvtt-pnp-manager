<script setup lang="ts">
import type { ImageEntry } from '~/types/entities'
import type { ApiImage, ApiMemory } from '~/types/api'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const campaign = useCampaignStore()
const id = route.params.id as string

const { data: image, error, refresh } = useAsyncData(
  `image-${id}`,
  async () => {
    await campaign.ensure()
    return api<ApiImage>(`/campaigns/${campaign.currentId}/images/${id}`)
  },
  { server: false },
)

const pushModes = [
  { v: 'empty_actor', l: 'Leerer Actor' },
  { v: 'journal', l: 'Journal / Handout' },
]

const form = reactive({ name: '', notes: '', pushAs: 'empty_actor', folderId: null as string | null })
const loaded = ref(false)
const imagePrompt = ref('')

watch(
  image,
  (im) => {
    if (!im || loaded.value) return
    form.name = im.name
    form.notes = im.notes
    form.pushAs = im.pushAs || 'empty_actor'
    form.folderId = im.folderId ?? null
    imagePrompt.value = im.name
    loaded.value = true
  },
  { immediate: true },
)

useSeoMeta({ title: () => form.name || t('nav.images') })

const folderType = computed(() => (form.pushAs === 'journal' ? 'JournalEntry' : 'Actor'))
const preview = computed<ImageEntry>(() => ({
  id,
  name: form.name || '—',
  image: image.value?.imageUrl || undefined,
  pushAs: form.pushAs,
  notes: form.notes,
}))

const saving = ref(false)
const saved = ref(false)
async function save() {
  saving.value = true
  saved.value = false
  try {
    await api(`/campaigns/${campaign.currentId}/images/${id}`, {
      method: 'PATCH',
      body: { name: form.name, notes: form.notes, pushAs: form.pushAs, folderId: form.folderId },
    })
    saved.value = true
    setTimeout(() => (saved.value = false), 2500)
  } finally {
    saving.value = false
  }
}

const confirming = ref(false)
async function remove() {
  if (!confirming.value) {
    confirming.value = true
    setTimeout(() => (confirming.value = false), 3000)
    return
  }
  await api(`/campaigns/${campaign.currentId}/images/${id}`, { method: 'DELETE' })
  await navigateTo('/images')
}

const pickerOpen = ref(false)
const pickedUrl = ref('')
watch(pickedUrl, () => refresh())

const syncing = ref(false)
const syncMsg = ref('')
async function syncToFoundry() {
  if (syncing.value || !campaign.currentId) return
  syncing.value = true
  syncMsg.value = ''
  try {
    await api(`/campaigns/${campaign.currentId}/images/${id}/sync`, { method: 'POST', body: {} })
    await refresh()
    syncMsg.value = 'synced'
  } catch (e: any) {
    const status = e?.response?.status ?? e?.statusCode
    syncMsg.value = status === 409 ? 'notConnected' : 'syncError'
  } finally {
    syncing.value = false
  }
}

const { data: imageMemories, refresh: refreshMemories } = useAsyncData(
  `image-memories-${id}`,
  async () => {
    await campaign.ensure()
    const rows = await api<ApiMemory[]>(`/campaigns/${campaign.currentId}/memories?subjectId=${id}`)
    return rows.map(toMemoryVM)
  },
  { server: false, default: () => [] },
)
async function ackMem(mid: string) {
  await api(`/campaigns/${campaign.currentId}/memories/${mid}`, { method: 'PATCH', body: { acknowledged: true } })
  await refreshMemories()
}
async function pinMem(m: { id: string; pinned?: boolean }) {
  await api(`/campaigns/${campaign.currentId}/memories/${m.id}`, { method: 'PATCH', body: { pinned: !m.pinned } })
  await refreshMemories()
}
async function delMem(mid: string) {
  await api(`/campaigns/${campaign.currentId}/memories/${mid}`, { method: 'DELETE' })
  await refreshMemories()
}
</script>

<template>
  <div class="editor">
    <NuxtLink to="/images" class="back"><Icon name="lucide:arrow-left" /> {{ t('image.back') }}</NuxtLink>

    <div v-if="error" class="notfound"><Icon name="lucide:ghost" /><p>{{ t('image.notFound') }}</p></div>

    <div v-else-if="loaded" class="cols">
      <div class="main">
        <AwPanel class="pad">
          <div class="grid">
            <div class="f wide">
              <label>{{ t('editor.name') }}</label>
              <input v-model="form.name" class="inp big" type="text" />
            </div>
            <div class="f">
              <label>{{ t('image.pushAs') }}</label>
              <select v-model="form.pushAs" class="inp">
                <option v-for="m in pushModes" :key="m.v" :value="m.v">{{ m.l }}</option>
              </select>
            </div>
            <div class="f">
              <label>{{ t('editor.folder') }}</label>
              <FolderPicker v-model="form.folderId" :type="folderType" />
            </div>
            <div class="f wide">
              <label>{{ t('image.notes') }}</label>
              <textarea v-model="form.notes" class="inp area" rows="3" :placeholder="t('image.notesPlaceholder')" />
            </div>
          </div>
        </AwPanel>

        <section class="mem-section">
          <div class="mh"><Icon name="lucide:sparkles" /><h3>{{ t('nav.memories') }}</h3></div>
          <MemoryComposer subject-type="image" :subject-id="id" :subject-label="form.name" :on-saved="refreshMemories" />
          <div v-if="imageMemories.length" class="mem-list">
            <MemoryCard v-for="m in imageMemories" :key="m.id" :memory="m" :on-ack="() => ackMem(m.id)" :on-pin="() => pinMem(m)" :on-delete="() => delMem(m.id)" />
          </div>
          <p v-else class="mem-empty">{{ t('memory.emptyImage') }}</p>
        </section>

        <EntityLinks subject-type="image" :subject-id="id" :subject-name="form.name" />

        <TagEditor subject-type="image" :subject-id="id" />
      </div>

      <aside class="rail">
        <div class="sticky">
          <span class="rl">{{ t('editor.preview') }}</span>
          <ImageCard :image="preview" />

          <div class="portrait">
            <span class="rl">{{ t('image.image') }}</span>
            <input v-model="imagePrompt" class="pinp" :placeholder="t('image.imagePrompt')" />
            <AwButton icon="lucide:images" variant="soft" @click="pickerOpen = true">
              {{ t('editor.chooseImage') }}
            </AwButton>
            <small class="phint">{{ t('editor.portraitHint') }}</small>
          </div>
          <ImagePicker
            v-model:open="pickerOpen"
            v-model:imageUrl="pickedUrl"
            subject-type="image"
            :subject-id="id"
            :prompt-hint="imagePrompt"
          />

          <div class="meta">
            <div class="row"><span>{{ t('editor.syncState') }}</span><b class="dirty">{{ image?.syncState || 'none' }}</b></div>
          </div>

          <div class="actions">
            <AwButton icon="lucide:save" variant="primary" @click="save">
              {{ saving ? t('editor.saving') : saved ? t('editor.saved') : t('editor.save') }}
            </AwButton>
            <AwButton icon="lucide:arrow-right-left" variant="ghost" @click="syncToFoundry">
              {{ syncing ? t('editor.syncing') : t('editor.sync') }}
            </AwButton>
            <p v-if="syncMsg" class="syncmsg" :class="syncMsg">
              <template v-if="syncMsg === 'synced'">{{ t('image.synced') }}</template>
              <template v-else-if="syncMsg === 'notConnected'">{{ t('editor.notConnected') }} <NuxtLink to="/settings/foundry">{{ t('nav.foundry') }} →</NuxtLink></template>
              <template v-else>{{ t('editor.syncError') }}</template>
            </p>
            <AwButton :icon="confirming ? 'lucide:trash-2' : 'lucide:trash'" variant="danger" @click="remove">
              {{ confirming ? t('editor.confirmDelete') : t('editor.delete') }}
            </AwButton>
          </div>
        </div>
      </aside>
    </div>

    <div v-else class="loading"><Icon name="lucide:loader-circle" class="spin" /></div>
  </div>
</template>

<style lang="scss" scoped>
.back {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin: 14px 0 22px;
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.08em;
  color: var(--ink-dim);
  text-decoration: none;
  transition: 0.2s;
  &:hover { color: var(--primary); }
  :deep(svg) { width: 15px; height: 15px; }
}

.cols { display: grid; grid-template-columns: 1fr 340px; gap: 24px; align-items: start; }
.main { display: flex; flex-direction: column; gap: 18px; }
.pad { padding: 22px; }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.f { display: flex; flex-direction: column; gap: 8px; }
.f.wide { grid-column: 1 / -1; }
.f label { font-family: var(--font-mono); font-size: 0.62rem; letter-spacing: 0.14em; text-transform: uppercase; color: var(--ink-faint); }
.inp {
  font-family: var(--font-body);
  font-size: 0.92rem;
  color: var(--ink);
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 11px;
  padding: 11px 13px;
  transition: 0.25s;
  &:focus { outline: 0; border-color: var(--primary); box-shadow: 0 0 0 3px rgba(70, 232, 255, 0.14), var(--glow-primary); }
  &.big { font-family: var(--font-display); font-size: 1.25rem; font-weight: 600; }
}
select.inp { cursor: pointer; appearance: none; }
.area { resize: vertical; line-height: 1.5; }

.rail .sticky { position: sticky; top: 90px; display: flex; flex-direction: column; gap: 16px; }
.rl { font-family: var(--font-mono); font-size: 0.6rem; letter-spacing: 0.2em; text-transform: uppercase; color: var(--ink-faint); }

.portrait {
  display: flex;
  flex-direction: column;
  gap: 9px;
  padding: 15px 16px;
  border-radius: 16px;
  border: 1px solid var(--line);
  background: var(--surface);
  .pinp {
    font-family: var(--font-body);
    font-size: 0.84rem;
    color: var(--ink);
    background: var(--surface-2);
    border: 1px solid var(--line-strong);
    border-radius: 10px;
    padding: 9px 12px;
    transition: 0.25s;
    &:focus { outline: 0; border-color: var(--secondary); box-shadow: var(--glow-secondary); }
  }
  .phint { font-size: 0.68rem; color: var(--ink-faint); line-height: 1.4; }
  :deep(.aw-btn) { width: 100%; justify-content: center; }
}

.meta { background: var(--surface); border: 1px solid var(--line); border-radius: 16px; padding: 6px 16px; }
.meta .row { display: flex; justify-content: space-between; padding: 11px 0; font-size: 0.82rem; span { color: var(--ink-faint); } b { font-family: var(--font-mono); font-size: 0.74rem; } .dirty { color: var(--gold); } }

.actions { display: flex; flex-direction: column; gap: 10px; }
.actions :deep(.aw-btn) { width: 100%; justify-content: center; }
.syncmsg { font-size: 0.74rem; font-family: var(--font-mono); text-align: center; line-height: 1.5; }
.syncmsg.synced { color: var(--emerald); }
.syncmsg.notConnected { color: var(--gold); }
.syncmsg.syncError { color: var(--ember); }
.syncmsg a { color: var(--primary); text-decoration: none; }

.mem-section { display: flex; flex-direction: column; gap: 12px; }
.mh { display: flex; align-items: center; gap: 10px; margin-top: 6px; :deep(svg) { width: 17px; height: 17px; color: var(--gold); } h3 { font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; } }
.mem-list { display: flex; flex-direction: column; gap: 10px; }
.mem-empty { font-size: 0.82rem; color: var(--ink-faint); font-family: var(--font-mono); padding: 4px 2px; }

.notfound, .loading { display: flex; flex-direction: column; align-items: center; gap: 14px; padding: 80px 20px; color: var(--ink-faint); :deep(svg) { width: 36px; height: 36px; } }
.spin { animation: spin 0.9s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 900px) {
  .cols { grid-template-columns: 1fr; }
  .rail .sticky { position: static; }
  .grid { grid-template-columns: 1fr; }
}
</style>
