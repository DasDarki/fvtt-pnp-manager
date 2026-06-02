<script setup lang="ts">
definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()

useSeoMeta({ title: () => t('foundry.title') })

interface Status {
  paired: boolean
  connected: boolean
  wsUrl: string
  world?: string
  version?: string
  lastSeenAt?: string | null
}

const status = ref<Status | null>(null)
const pairing = ref<{ token: string; wsUrl: string } | null>(null)
const busy = ref(false)
const copied = ref('')

async function loadStatus() {
  await campaign.ensure()
  status.value = await api<Status>(`/campaigns/${campaign.currentId}/foundry`)
}

async function pair() {
  busy.value = true
  try {
    pairing.value = await api(`/campaigns/${campaign.currentId}/foundry/pair`, { method: 'POST', body: {} })
    await loadStatus()
  } finally {
    busy.value = false
  }
}

async function copy(value: string, key: string) {
  try {
    await navigator.clipboard.writeText(value)
    copied.value = key
    setTimeout(() => (copied.value = ''), 1500)
  } catch {
    /* ignore */
  }
}

interface Folder {
  id: string
  name: string
  foundryType: string
  parentId: string | null
  foundryFolderId: string
}

const folders = ref<Folder[]>([])
const discovering = ref(false)

const typeIcons: Record<string, string> = {
  Actor: 'lucide:users',
  Item: 'lucide:gem',
  Scene: 'lucide:castle',
  JournalEntry: 'lucide:scroll-text',
}
function typeIcon(type: string) {
  return typeIcons[type] || 'lucide:folder'
}

const folderGroups = computed(() => {
  const byType: Record<string, Folder[]> = {}
  for (const f of folders.value) (byType[f.foundryType || '—'] ??= []).push(f)
  return Object.entries(byType).map(([type, list]) => {
    const ids = new Set(list.map((f) => f.id))
    const children: Record<string, Folder[]> = {}
    for (const f of list) {
      const key = f.parentId && ids.has(f.parentId) ? f.parentId : '__root'
      ;(children[key] ??= []).push(f)
    }
    const rows: { id: string; name: string; depth: number }[] = []
    const walk = (pid: string, depth: number) => {
      for (const f of (children[pid] || []).slice().sort((a, b) => a.name.localeCompare(b.name))) {
        rows.push({ id: f.id, name: f.name, depth })
        walk(f.id, depth + 1)
      }
    }
    walk('__root', 0)
    return { type, rows }
  })
})

async function loadFolders() {
  await campaign.ensure()
  folders.value = await api<Folder[]>(`/campaigns/${campaign.currentId}/folders`)
}

async function discover() {
  if (discovering.value) return
  discovering.value = true
  try {
    const res = await api<{ folders: Folder[] }>(`/campaigns/${campaign.currentId}/foundry/discover`, {
      method: 'POST',
      body: {},
    })
    folders.value = res.folders
  } finally {
    discovering.value = false
  }
}

let timer: ReturnType<typeof setInterval> | null = null
onMounted(() => {
  loadStatus()
  loadFolders()
  timer = setInterval(loadStatus, 5000)
})
onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="foundry">
    <header class="head">
      <span class="ic"><Icon name="lucide:plug-zap" /></span>
      <div>
        <h1>{{ t('foundry.title') }}</h1>
        <p>{{ t('foundry.subtitle') }}</p>
      </div>
    </header>

    <div class="grid">
      <AwPanel class="card">
        <div class="ch"><Icon name="lucide:activity" /><h2>{{ t('foundry.statusTitle') }}</h2></div>

        <div class="state" :class="{ on: status?.connected }">
          <span class="led" />
          <b v-if="status?.connected">{{ t('foundry.connected') }}</b>
          <b v-else-if="status?.paired">{{ t('foundry.disconnected') }}</b>
          <b v-else>{{ t('foundry.notPaired') }}</b>
        </div>

        <div v-if="status?.connected" class="rows">
          <div class="row"><span>{{ t('foundry.world') }}</span><b>{{ status.world || '—' }}</b></div>
          <div class="row"><span>Foundry</span><b>v{{ status.version || '—' }}</b></div>
        </div>

        <AwButton icon="lucide:refresh-cw" variant="ghost" @click="loadStatus">{{ t('foundry.refresh') }}</AwButton>
      </AwPanel>

      <AwPanel class="card">
        <div class="ch"><Icon name="lucide:key-round" /><h2>{{ t('foundry.pairTitle') }}</h2></div>
        <p class="hint">{{ t('foundry.pairHint') }}</p>

        <AwButton icon="lucide:sparkles" variant="primary" :disabled="busy" @click="pair">
          {{ status?.paired ? t('foundry.regenerate') : t('foundry.generate') }}
        </AwButton>

        <div v-if="pairing" class="creds">
          <p class="once">{{ t('foundry.tokenOnce') }}</p>
          <label class="cred">
            <span>{{ t('foundry.relayUrl') }}</span>
            <button class="field" @click="copy(pairing.wsUrl, 'url')">
              <code>{{ pairing.wsUrl }}</code>
              <Icon :name="copied === 'url' ? 'lucide:check' : 'lucide:copy'" />
            </button>
          </label>
          <label class="cred">
            <span>{{ t('foundry.token') }}</span>
            <button class="field" @click="copy(pairing.token, 'tok')">
              <code>{{ pairing.token }}</code>
              <Icon :name="copied === 'tok' ? 'lucide:check' : 'lucide:copy'" />
            </button>
          </label>
        </div>

        <ol class="steps">
          <li>{{ t('foundry.step1') }}</li>
          <li>{{ t('foundry.step2') }}</li>
          <li>{{ t('foundry.step3') }}</li>
        </ol>
      </AwPanel>
    </div>

    <AwPanel class="card folders">
      <div class="ch">
        <Icon name="lucide:folder-tree" />
        <h2>{{ t('foundry.foldersTitle') }}</h2>
        <AwButton
          class="discover-btn"
          icon="lucide:folder-sync"
          variant="ghost"
          :disabled="!status?.connected || discovering"
          @click="discover"
        >
          {{ discovering ? t('foundry.discovering') : t('foundry.discover') }}
        </AwButton>
      </div>
      <p class="hint">{{ t('foundry.foldersHint') }}</p>

      <div v-if="folderGroups.length" class="ftree">
        <div v-for="g in folderGroups" :key="g.type" class="fgroup">
          <div class="ftype"><Icon :name="typeIcon(g.type)" /> {{ g.type }} <span>{{ g.rows.length }}</span></div>
          <div v-for="r in g.rows" :key="r.id" class="frow" :style="{ paddingLeft: `${10 + r.depth * 20}px` }">
            <Icon name="lucide:folder" /> {{ r.name }}
          </div>
        </div>
      </div>
      <p v-else class="empty-f">{{ t('foundry.noFolders') }}</p>
    </AwPanel>
  </div>
</template>

<style lang="scss" scoped>
.head {
  display: flex;
  align-items: center;
  gap: 16px;
  margin: 14px 0 26px;

  .ic {
    width: 50px;
    height: 50px;
    border-radius: 15px;
    display: grid;
    place-items: center;
    flex: none;
    background: color-mix(in srgb, var(--emerald) 16%, transparent);
    color: var(--emerald);
    border: 1px solid var(--line);
    box-shadow: 0 0 22px -6px var(--emerald);
    :deep(svg) { width: 24px; height: 24px; }
  }
  h1 { font-family: var(--font-display); font-weight: 600; font-size: clamp(1.6rem, 3vw, 2.2rem); letter-spacing: 0.02em; text-shadow: var(--glow-text); }
  p { color: var(--ink-faint); font-size: 0.84rem; margin-top: 4px; }
}

.grid { display: grid; grid-template-columns: 1fr 1.2fr; gap: 22px; align-items: start; }
.card { padding: 24px; display: flex; flex-direction: column; gap: 16px; }
.ch {
  display: flex;
  align-items: center;
  gap: 10px;
  :deep(svg) { width: 17px; height: 17px; color: var(--primary); }
  h2 { font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; }
}
.hint { font-size: 0.82rem; color: var(--ink-faint); }

.state {
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 14px 16px;
  border-radius: 13px;
  border: 1px solid var(--line);
  background: var(--surface-2);

  .led { width: 11px; height: 11px; border-radius: 50%; background: var(--ink-faint); }
  b { font-family: var(--font-display); font-weight: 600; font-size: 1rem; }
  &.on {
    border-color: rgba(55, 232, 164, 0.4);
    background: rgba(55, 232, 164, 0.08);
    .led { background: var(--emerald); box-shadow: 0 0 10px var(--emerald); animation: blink 2s ease-in-out infinite; }
    b { color: var(--emerald); }
  }
}
@keyframes blink { 50% { opacity: 0.4; } }

.rows .row {
  display: flex;
  justify-content: space-between;
  padding: 10px 2px;
  border-bottom: 1px solid var(--line);
  font-size: 0.84rem;
  &:last-child { border-bottom: 0; }
  span { color: var(--ink-faint); }
  b { font-family: var(--font-mono); font-size: 0.78rem; }
}

.creds {
  display: flex;
  flex-direction: column;
  gap: 11px;
  padding: 16px;
  border-radius: 14px;
  border: 1px dashed var(--line-strong);
  background: var(--surface);

  .once { font-family: var(--font-mono); font-size: 0.66rem; color: var(--gold); letter-spacing: 0.04em; }
  .cred { display: flex; flex-direction: column; gap: 6px; }
  .cred > span { font-family: var(--font-mono); font-size: 0.58rem; letter-spacing: 0.14em; text-transform: uppercase; color: var(--ink-faint); }
  .field {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 13px;
    border-radius: 10px;
    border: 1px solid var(--line-strong);
    background: var(--surface-2);
    cursor: pointer;
    transition: 0.2s;
    text-align: left;

    &:hover { border-color: var(--primary); }
    code { flex: 1; font-family: var(--font-mono); font-size: 0.78rem; color: var(--ink); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
    :deep(svg) { width: 15px; height: 15px; color: var(--primary); flex: none; }
  }
}

.steps {
  margin: 0;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 0.82rem;
  color: var(--ink-dim);
  li::marker { color: var(--primary); font-family: var(--font-mono); }
}

.folders { margin-top: 22px; }
.ch .discover-btn { margin-left: auto; }
.ftree { display: flex; flex-direction: column; gap: 18px; margin-top: 6px; }
.ftype {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 0.6rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--ink-faint);
  padding-bottom: 8px;
  border-bottom: 1px solid var(--line);
  margin-bottom: 4px;
  :deep(svg) { width: 14px; height: 14px; color: var(--primary); }
  span { margin-left: auto; color: var(--ink-dim); }
}
.frow {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
  font-size: 0.86rem;
  color: var(--ink-dim);
  :deep(svg) { width: 14px; height: 14px; color: var(--ink-faint); flex: none; }
}
.empty-f { font-size: 0.82rem; color: var(--ink-faint); font-family: var(--font-mono); margin-top: 4px; }

@media (max-width: 900px) {
  .grid { grid-template-columns: 1fr; }
}
</style>
