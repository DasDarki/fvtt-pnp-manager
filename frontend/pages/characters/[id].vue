<script setup lang="ts">
import type { Character, CharacterStatus, CharacterType } from '~/types/entities'
import type { ApiCharacter, ApiMemory } from '~/types/api'
import type { SystemData } from '~/types/ruleset'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const campaign = useCampaignStore()
const id = route.params.id as string

const { data: char, pending, error, refresh } = useAsyncData(
  `character-${id}`,
  async () => {
    await campaign.ensure()
    return api<ApiCharacter>(`/campaigns/${campaign.currentId}/characters/${id}`)
  },
  { server: false },
)

const adapter = computed(() => getAdapter(campaign.current?.ruleset))

const form = reactive({
  name: '',
  status: 'alive' as CharacterStatus,
  characterType: 'npc' as CharacterType,
  folderId: null as string | null,
  system: {} as SystemData,
})
const loaded = ref(false)

const portraitPrompt = ref('')
const imageAlign = ref('center')

watch(
  char,
  (c) => {
    if (!c || loaded.value) return
    form.name = c.name
    form.status = (c.status as CharacterStatus) || 'alive'
    form.characterType = (c.characterType as CharacterType) || 'npc'
    form.folderId = c.folderId ?? null
    form.system = reactive(mergeSystemData(adapter.value.schema, c.systemData))
    imageAlign.value = (c.systemData as any)?.imageAlign || 'center'
    portraitPrompt.value = `Porträt von ${c.name}`
    loaded.value = true
  },
  { immediate: true },
)

useSeoMeta({ title: () => form.name || t('nav.characters') })

const { data: charMemories, refresh: refreshMemories } = useAsyncData(
  `char-memories-${id}`,
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

const statuses: CharacterStatus[] = ['alive', 'dead', 'unknown', 'hunted']
const types: CharacterType[] = ['pc', 'npc', 'ally', 'foe', 'neutral']

const preview = computed<Character>(() => ({
  id,
  name: form.name || '—',
  subtitle: adapter.value.subtitle(form.system),
  type: form.characterType,
  status: form.status,
  initial: (form.name || '?').charAt(0).toUpperCase(),
  ring: form.system.ring || 'conic-gradient(from 140deg,var(--primary),var(--secondary),var(--magenta),var(--primary))',
  image: char.value?.imageUrl || undefined,
  imageAlign: imageAlign.value,
  stats: adapter.value.cardStats(form.system),
  hpPercent: adapter.value.hpPercent(form.system),
  critical: !!form.system.critical,
}))

const syncing = ref(false)
const syncMsg = ref('')
async function syncToFoundry() {
  if (syncing.value || !campaign.currentId) return
  syncing.value = true
  syncMsg.value = ''
  try {
    await api(`/campaigns/${campaign.currentId}/characters/${id}/sync`, { method: 'POST', body: {} })
    await refresh()
    syncMsg.value = 'synced'
  } catch (e: any) {
    const status = e?.response?.status ?? e?.statusCode
    syncMsg.value = status === 409 ? 'notConnected' : 'syncError'
  } finally {
    syncing.value = false
  }
}

const pickerOpen = ref(false)
const pickedUrl = ref('')
watch(pickedUrl, () => refresh())

const saving = ref(false)
const saved = ref(false)
async function save() {
  saving.value = true
  saved.value = false
  try {
    await api(`/campaigns/${campaign.currentId}/characters/${id}`, {
      method: 'PATCH',
      body: {
        name: form.name,
        status: form.status,
        characterType: form.characterType,
        folderId: form.folderId,
        systemData: { ...form.system, imageAlign: imageAlign.value },
      },
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
  await api(`/campaigns/${campaign.currentId}/characters/${id}`, { method: 'DELETE' })
  await navigateTo('/characters')
}
</script>

<template>
  <div class="editor">
    <NuxtLink to="/characters" class="back">
      <Icon name="lucide:arrow-left" /> {{ t('editor.back') }}
    </NuxtLink>

    <div v-if="error" class="notfound">
      <Icon name="lucide:ghost" />
      <p>{{ t('editor.notFound') }}</p>
      <NuxtLink to="/characters" class="link">{{ t('editor.back') }}</NuxtLink>
    </div>

    <div v-else-if="loaded" class="cols">
      <div class="main">
        <AwPanel class="ident">
          <div class="grid">
            <div class="sf wide">
              <label class="lbl">{{ t('editor.name') }}</label>
              <input v-model="form.name" class="inp big" type="text" />
            </div>
            <div class="sf">
              <label class="lbl">{{ t('editor.type') }}</label>
              <select v-model="form.characterType" class="inp">
                <option v-for="ty in types" :key="ty" :value="ty">{{ t(`charType.${ty}`) }}</option>
              </select>
            </div>
            <div class="sf">
              <label class="lbl">{{ t('editor.status') }}</label>
              <div class="statuses">
                <button
                  v-for="s in statuses"
                  :key="s"
                  class="st"
                  :class="[s, { on: form.status === s }]"
                  @click="form.status = s"
                >
                  {{ t(`status.${s}`) }}
                </button>
              </div>
            </div>
            <div class="sf wide">
              <label class="lbl">{{ t('editor.folder') }}</label>
              <FolderPicker v-model="form.folderId" type="Actor" />
            </div>
          </div>
        </AwPanel>

        <SchemaForm :schema="adapter.schema" :data="form.system" />

        <section class="mem-section">
          <div class="mh">
            <Icon name="lucide:sparkles" />
            <h3>{{ t('nav.memories') }}</h3>
          </div>
          <MemoryComposer
            subject-type="character"
            :subject-id="id"
            :subject-label="form.name"
            :on-saved="refreshMemories"
          />
          <div v-if="charMemories.length" class="mem-list">
            <MemoryCard
              v-for="m in charMemories"
              :key="m.id"
              :memory="m"
              :on-ack="() => ackMem(m.id)"
              :on-pin="() => pinMem(m)"
              :on-delete="() => delMem(m.id)"
            />
          </div>
          <p v-else class="mem-empty">{{ t('memory.emptyChar') }}</p>
        </section>

        <EntityLinks subject-type="character" :subject-id="id" :subject-name="form.name" />

        <TagEditor subject-type="character" :subject-id="id" />
      </div>

      <aside class="rail">
        <div class="sticky">
          <span class="rl">{{ t('editor.preview') }}</span>
          <CharacterCard :character="preview" />

          <div class="portrait">
            <span class="rl">{{ t('editor.portrait') }}</span>
            <input v-model="portraitPrompt" class="pinp" :placeholder="t('editor.portraitPrompt')" />
            <AwButton icon="lucide:images" variant="soft" @click="pickerOpen = true">
              {{ t('editor.chooseImage') }}
            </AwButton>
            <ImageAlignToggle v-if="preview.image" v-model="imageAlign" @update:model-value="save" />
            <small class="phint">{{ t('editor.portraitHint') }}</small>
          </div>
          <ImagePicker
            v-model:open="pickerOpen"
            v-model:imageUrl="pickedUrl"
            subject-type="character"
            :subject-id="id"
            :prompt-hint="portraitPrompt"
          />

          <div class="meta">
            <div class="row">
              <span>{{ t('editor.ruleset') }}</span>
              <b>{{ adapter.label }}</b>
            </div>
            <div class="row">
              <span>{{ t('editor.syncState') }}</span>
              <b class="dirty">{{ char?.syncState || 'none' }}</b>
            </div>
          </div>

          <div class="actions">
            <AwButton icon="lucide:save" variant="primary" @click="save">
              {{ saving ? t('editor.saving') : saved ? t('editor.saved') : t('editor.save') }}
            </AwButton>
            <AwButton icon="lucide:arrow-right-left" variant="ghost" @click="syncToFoundry">
              {{ syncing ? t('editor.syncing') : t('editor.sync') }}
            </AwButton>
            <p v-if="syncMsg" class="syncmsg" :class="syncMsg">
              <template v-if="syncMsg === 'synced'">{{ t('editor.synced') }}</template>
              <template v-else-if="syncMsg === 'notConnected'">
                {{ t('editor.notConnected') }} <NuxtLink to="/settings/foundry">{{ t('nav.foundry') }} →</NuxtLink>
              </template>
              <template v-else>{{ t('editor.syncError') }}</template>
            </p>
            <AwButton :icon="confirming ? 'lucide:trash-2' : 'lucide:trash'" variant="danger" @click="remove">
              {{ confirming ? t('editor.confirmDelete') : t('editor.delete') }}
            </AwButton>
          </div>
        </div>
      </aside>
    </div>

    <div v-else class="loading">
      <Icon name="lucide:loader-circle" class="spin" />
    </div>
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

.cols {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 24px;
  align-items: start;
}

.main { display: flex; flex-direction: column; gap: 18px; }

.mem-section { display: flex; flex-direction: column; gap: 12px; }
.mh {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 6px;
  :deep(svg) { width: 17px; height: 17px; color: var(--gold); }
  h3 { font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; letter-spacing: 0.02em; }
}
.mem-list { display: flex; flex-direction: column; gap: 10px; }
.mem-empty { font-size: 0.82rem; color: var(--ink-faint); font-family: var(--font-mono); padding: 4px 2px; }

.ident { padding: 22px; }
.ident .grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
.sf { display: flex; flex-direction: column; gap: 8px; }
.sf.wide { grid-column: 1 / -1; }
.lbl {
  font-family: var(--font-mono);
  font-size: 0.62rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--ink-faint);
}
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

.statuses { display: flex; flex-wrap: wrap; gap: 7px; }
.st {
  font-family: var(--font-mono);
  font-size: 0.64rem;
  letter-spacing: 0.06em;
  padding: 8px 12px;
  border-radius: 9px;
  border: 1px solid var(--line);
  background: var(--surface);
  color: var(--ink-dim);
  cursor: pointer;
  transition: 0.2s;

  &:hover { border-color: var(--line-strong); color: var(--ink); }
  &.on { color: #06040c; border-color: transparent; background: var(--grad-arcane); box-shadow: var(--glow-primary); }
  &.on.dead { background: var(--ember); color: #1a0a06; box-shadow: 0 0 18px rgba(255, 106, 85, 0.4); }
  &.on.hunted { background: var(--gold); color: #1a1206; box-shadow: 0 0 18px rgba(255, 194, 77, 0.4); }
}
:global(html[data-theme='light']) .st.on { color: #fff; }

.rail .sticky {
  position: sticky;
  top: 90px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.rl {
  font-family: var(--font-mono);
  font-size: 0.6rem;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: var(--ink-faint);
}

.meta {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 16px;
  padding: 6px 16px;

  .row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 11px 0;
    border-bottom: 1px solid var(--line);
    font-size: 0.82rem;

    &:last-child { border-bottom: 0; }
    span { color: var(--ink-faint); }
    b { font-family: var(--font-mono); font-size: 0.74rem; }
    .dirty { color: var(--gold); }
  }
}

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

.actions { display: flex; flex-direction: column; gap: 10px; }
.actions :deep(.aw-btn) { width: 100%; justify-content: center; }
.syncmsg { font-size: 0.74rem; font-family: var(--font-mono); text-align: center; line-height: 1.5; }
.syncmsg.synced { color: var(--emerald); }
.syncmsg.notConnected { color: var(--gold); }
.syncmsg.syncError { color: var(--ember); }
.syncmsg a { color: var(--primary); text-decoration: none; }

.notfound,
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding: 80px 20px;
  color: var(--ink-faint);

  :deep(svg) { width: 36px; height: 36px; }
  .link { color: var(--primary); font-family: var(--font-mono); font-size: 0.78rem; text-decoration: none; }
}
.spin { animation: spin 0.9s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 900px) {
  .cols { grid-template-columns: 1fr; }
  .rail .sticky { position: static; }
}
</style>
