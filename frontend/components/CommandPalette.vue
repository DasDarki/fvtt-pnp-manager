<script setup lang="ts">
import type { ApiCharacter, ApiItem, ApiScene } from '~/types/api'

interface Entry {
  label: string
  sub: string
  icon: string
  group: string
  run: () => void | Promise<void>
}

const { t } = useI18n()
const ui = useUiStore()
const api = useApi()
const campaign = useCampaignStore()

const query = ref('')
const selected = ref(0)
const inputEl = ref<HTMLInputElement | null>(null)
const listEl = ref<HTMLElement | null>(null)

const characters = ref<{ id: string; name: string; sub: string }[]>([])
const items = ref<{ id: string; name: string; sub: string }[]>([])
const scenes = ref<{ id: string; name: string }[]>([])
const loaded = ref(false)
const busy = ref(false)

function go(path: string) {
  ui.closePalette()
  navigateTo(path)
}

async function createCharacter() {
  if (busy.value || !campaign.currentId) return
  busy.value = true
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
    ui.closePalette()
    await navigateTo(`/characters/${ch.id}`)
  } finally {
    busy.value = false
  }
}

const actionDefs = computed<Entry[]>(() => [
  { label: t('palette.action.new-character.title'), sub: t('palette.action.new-character.sub'), icon: 'lucide:user-plus', group: t('palette.actions'), run: createCharacter },
  { label: t('palette.action.generate-image.title'), sub: t('palette.action.generate-image.sub'), icon: 'lucide:wand-sparkles', group: t('palette.actions'), run: () => go('/dalle') },
  { label: t('nav.dashboard'), sub: t('meta.tagline'), icon: 'lucide:layout-dashboard', group: t('palette.actions'), run: () => go('/dashboard') },
  { label: t('nav.foundry'), sub: t('foundry.subtitle'), icon: 'lucide:plug-zap', group: t('palette.actions'), run: () => go('/settings/foundry') },
])

const flat = computed<Entry[]>(() => {
  const q = query.value.toLowerCase().trim()
  const match = (s: string) => !q || s.toLowerCase().includes(q)
  const out: Entry[] = []

  for (const a of actionDefs.value) if (match(a.label)) out.push(a)

  for (const c of characters.value.filter((x) => match(x.name)).slice(0, 6)) {
    out.push({ label: c.name, sub: c.sub, icon: 'lucide:user', group: t('nav.characters'), run: () => go(`/characters/${c.id}`) })
  }
  for (const i of items.value.filter((x) => match(x.name)).slice(0, 5)) {
    out.push({ label: i.name, sub: i.sub, icon: 'lucide:gem', group: t('nav.items'), run: () => go('/items') })
  }
  for (const s of scenes.value.filter((x) => match(x.name)).slice(0, 5)) {
    out.push({ label: s.name, sub: t('nav.scenes'), icon: 'lucide:castle', group: t('nav.scenes'), run: () => go('/scenes') })
  }
  return out
})

const grouped = computed(() => {
  const map = new Map<string, { entry: Entry; index: number }[]>()
  flat.value.forEach((entry, index) => {
    if (!map.has(entry.group)) map.set(entry.group, [])
    map.get(entry.group)!.push({ entry, index })
  })
  return [...map.entries()].map(([title, entries]) => ({ title, entries }))
})

watch(query, () => (selected.value = 0))
watch(flat, () => {
  if (selected.value >= flat.value.length) selected.value = Math.max(0, flat.value.length - 1)
})

async function loadData() {
  await campaign.ensure()
  const cid = campaign.currentId
  const [c, i, s] = await Promise.all([
    api<ApiCharacter[]>(`/campaigns/${cid}/characters`),
    api<ApiItem[]>(`/campaigns/${cid}/items`),
    api<ApiScene[]>(`/campaigns/${cid}/scenes`),
  ])
  characters.value = c.map((x) => ({ id: x.id, name: x.name, sub: x.systemData?.subtitle || x.status }))
  items.value = i.map((x) => ({ id: x.id, name: x.name, sub: x.itemType || x.rarity }))
  scenes.value = s.map((x) => ({ id: x.id, name: x.name }))
  loaded.value = true
}

watch(
  () => ui.paletteOpen,
  (open) => {
    if (open) {
      query.value = ''
      selected.value = 0
      nextTick(() => inputEl.value?.focus())
      loadData()
    }
  },
)

function onKey(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    ui.togglePalette()
    return
  }
  if (!ui.paletteOpen) return
  if (e.key === 'Escape') {
    ui.closePalette()
  } else if (e.key === 'ArrowDown') {
    e.preventDefault()
    selected.value = Math.min(selected.value + 1, flat.value.length - 1)
    scrollSelected()
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    selected.value = Math.max(selected.value - 1, 0)
    scrollSelected()
  } else if (e.key === 'Enter') {
    e.preventDefault()
    flat.value[selected.value]?.run()
  }
}

function scrollSelected() {
  nextTick(() => {
    listEl.value?.querySelector<HTMLElement>('.res.sel')?.scrollIntoView({ block: 'nearest' })
  })
}

onMounted(() => window.addEventListener('keydown', onKey))
onBeforeUnmount(() => window.removeEventListener('keydown', onKey))
</script>

<template>
  <Teleport to="body">
    <Transition name="pal">
      <div v-if="ui.paletteOpen" class="palette" @click.self="ui.closePalette()">
        <div class="box">
          <div class="pin">
            <Icon name="lucide:search" />
            <input ref="inputEl" v-model="query" :placeholder="t('palette.placeholder')" />
            <Icon v-if="busy" name="lucide:loader-circle" class="spin" />
          </div>

          <div ref="listEl" class="list">
            <template v-for="g in grouped" :key="g.title">
              <div class="grp">{{ g.title }}</div>
              <button
                v-for="row in g.entries"
                :key="row.index"
                class="res"
                :class="{ sel: row.index === selected }"
                @mouseenter="selected = row.index"
                @click="row.entry.run()"
              >
                <span class="ri"><Icon :name="row.entry.icon" /></span>
                <span class="rt">
                  <b>{{ row.entry.label }}</b>
                  <small>{{ row.entry.sub }}</small>
                </span>
                <span v-if="row.index === selected" class="rk">↵</span>
              </button>
            </template>

            <div v-if="flat.length === 0" class="none">{{ t('palette.empty') }}</div>
          </div>

          <div class="foot">
            <span><span class="kbd">↑↓</span> {{ t('palette.navigate') }}</span>
            <span><span class="kbd">↵</span> {{ t('palette.open') }}</span>
            <span><span class="kbd">Esc</span> {{ t('palette.close') }}</span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.palette {
  position: fixed;
  inset: 0;
  z-index: 80;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 13vh;
  background: rgba(5, 4, 10, 0.6);
  backdrop-filter: blur(8px);
}

.box {
  width: min(620px, 92vw);
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 22px;
  box-shadow: 0 40px 90px -30px #000, var(--glow-secondary);
  overflow: hidden;
}

.pin {
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 18px 20px;
  border-bottom: 1px solid var(--line);

  :deep(svg) { width: 20px; height: 20px; color: var(--primary); }
  input { flex: 1; background: transparent; border: 0; outline: 0; color: var(--ink); font-size: 1.05rem; font-family: var(--font-body); }
  .spin { animation: spin 0.9s linear infinite; }
}

.list { max-height: 52vh; overflow-y: auto; padding-bottom: 6px; }

.grp {
  font-family: var(--font-mono);
  font-size: 0.58rem;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: var(--ink-faint);
  padding: 14px 20px 6px;
}

.res {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 11px 20px;
  cursor: pointer;
  transition: background 0.12s;
  border: 0;
  background: transparent;
  text-align: left;

  &.sel { background: rgba(70, 232, 255, 0.1); }
  .ri {
    width: 32px;
    height: 32px;
    border-radius: 9px;
    display: grid;
    place-items: center;
    flex: none;
    background: var(--surface);
    border: 1px solid var(--line);
    color: var(--primary);
    :deep(svg) { width: 16px; height: 16px; }
  }
  .rt { flex: 1; min-width: 0; }
  .rt b { font-size: 0.9rem; font-weight: 600; color: var(--ink); display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .rt small { display: block; font-size: 0.72rem; color: var(--ink-faint); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .rk { font-family: var(--font-mono); font-size: 0.7rem; color: var(--primary); }
}

.none { padding: 26px 20px; text-align: center; color: var(--ink-faint); font-family: var(--font-mono); font-size: 0.8rem; }

.foot {
  display: flex;
  gap: 18px;
  padding: 13px 20px;
  border-top: 1px solid var(--line);
  font-family: var(--font-mono);
  font-size: 0.6rem;
  color: var(--ink-faint);
  span { display: flex; align-items: center; gap: 6px; }
}

.kbd {
  font-family: var(--font-mono);
  font-size: 0.6rem;
  color: var(--ink-dim);
  border: 1px solid var(--line-strong);
  border-radius: 6px;
  padding: 3px 7px;
  background: var(--surface);
}

.pal-enter-active,
.pal-leave-active { transition: opacity 0.2s ease; }
.pal-enter-from,
.pal-leave-to { opacity: 0; }
.pal-enter-active .box { animation: pop 0.25s cubic-bezier(0.2, 0.9, 0.3, 1.2); }

@keyframes pop { from { transform: translateY(-14px) scale(0.97); opacity: 0; } }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
