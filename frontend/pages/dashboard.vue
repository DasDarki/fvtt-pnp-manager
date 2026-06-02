<script setup lang="ts">
import type { StatTile } from '~/types/entities'
import type { ApiCharacter, ApiItem, ApiMemory, ApiScene } from '~/types/api'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const ui = useUiStore()
const api = useApi()
const campaign = useCampaignStore()

useSeoMeta({ title: () => t('nav.dashboard') })

const { data, refresh } = useAsyncData(
  'dashboard',
  async () => {
    await campaign.ensure()
    const cid = campaign.currentId
    const [chars, scenes, items, memories, alerts] = await Promise.all([
      api<ApiCharacter[]>(`/campaigns/${cid}/characters`),
      api<ApiScene[]>(`/campaigns/${cid}/scenes`),
      api<ApiItem[]>(`/campaigns/${cid}/items`),
      api<ApiMemory[]>(`/campaigns/${cid}/memories`),
      api<ApiMemory[]>(`/campaigns/${cid}/alerts`),
    ])
    return {
      characters: chars.map(toCharacterVM),
      scenes: scenes.map(toSceneVM),
      itemCount: items.length,
      memoryCount: memories.length,
      memories: memories.map(toMemoryVM),
      alerts: alerts.map(toMemoryVM),
    }
  },
  {
    server: false,
    default: () => ({ characters: [], scenes: [], itemCount: 0, memoryCount: 0, memories: [], alerts: [] }),
  },
)

const recent = computed(() => data.value.characters.slice(0, 6))
const feed = computed(() => data.value.memories.slice(0, 4))
const featured = computed(() => data.value.scenes[0])
const firstAlert = computed(() => data.value.alerts[0])

const stats = computed<StatTile[]>(() => [
  { key: 'characters', icon: 'lucide:users', value: String(data.value.characters.length), accent: 'primary' },
  { key: 'scenes', icon: 'lucide:castle', value: String(data.value.scenes.length), accent: 'secondary' },
  { key: 'items', icon: 'lucide:gem', value: String(data.value.itemCount), accent: 'gold' },
  { key: 'memories', icon: 'lucide:sparkles', value: String(data.value.memoryCount), accent: 'magenta' },
])

async function ackAlert(id: string) {
  await api(`/campaigns/${campaign.currentId}/memories/${id}`, {
    method: 'PATCH',
    body: { acknowledged: true },
  })
  await refresh()
}

const quick = [
  { key: 'character', icon: 'lucide:user-plus', accent: 'primary' },
  { key: 'item', icon: 'lucide:gem', accent: 'gold' },
  { key: 'image', icon: 'lucide:wand-sparkles', accent: 'secondary' },
  { key: 'sync', icon: 'lucide:arrow-right-left', accent: 'emerald' },
]
</script>

<template>
  <div class="dash">
    <div class="hello">
      <div>
        <span class="eyebrow">{{ t('dash.welcome') }}</span>
        <h1>{{ t('dash.campaign') }} <em>{{ campaign.current?.name || '—' }}</em></h1>
        <p>{{ t('dash.subtitle') }}</p>
      </div>
      <AwButton icon="lucide:wand-sparkles" variant="primary" class="hide-m" @click="ui.openPalette()">
        {{ t('actions.generate') }}
      </AwButton>
    </div>

    <CriticalAlert
      v-if="firstAlert"
      :title="firstAlert.title"
      :body="firstAlert.body"
      :subject-label="firstAlert.subjectLabel"
      :more-count="data.alerts.length - 1"
      :on-ack="() => ackAlert(firstAlert.id)"
    />

    <div class="stats">
      <StatTile v-for="s in stats" :key="s.key" :stat="s" />
    </div>

    <div class="cols">
      <div>
        <div class="block">
          <div class="block-head">
            <h2>{{ t('dash.recent') }}</h2>
            <div class="rule" />
            <NuxtLink to="/characters">{{ t('dash.viewAll') }} →</NuxtLink>
          </div>
          <div class="charscroll">
            <NuxtLink v-for="c in recent" :key="c.id" :to="`/characters/${c.id}`" class="card-link">
              <CharacterCard :character="c" />
            </NuxtLink>
          </div>
        </div>

        <div v-if="featured" class="block">
          <div class="block-head">
            <h2>{{ t('dash.preparedScene') }}</h2>
            <div class="rule" />
            <NuxtLink to="/scenes">{{ t('dash.openFoundry') }} →</NuxtLink>
          </div>
          <SceneCard :scene="featured" />
        </div>
      </div>

      <div>
        <div class="block">
          <div class="block-head">
            <h2>{{ t('nav.memories') }}</h2>
            <div class="rule" />
            <a href="#">{{ t('dash.all') }} →</a>
          </div>
          <AwPanel>
            <div class="feed">
              <MemoryItem v-for="(m, i) in feed" :key="m.id" :memory="m" :last="i === feed.length - 1" />
            </div>
          </AwPanel>
        </div>

        <div class="block">
          <div class="block-head">
            <h2>{{ t('dash.quick') }}</h2>
            <div class="rule" />
          </div>
          <AwPanel>
            <div class="quick">
              <button
                v-for="q in quick"
                :key="q.key"
                class="qa"
                :style="{ '--c': `var(--${q.accent})` }"
                @click="ui.openPalette()"
              >
                <span class="qi"><Icon :name="q.icon" /></span>
                <span class="qt">
                  <b>{{ t(`dash.quickItem.${q.key}.title`) }}</b>
                  <small>{{ t(`dash.quickItem.${q.key}.sub`) }}</small>
                </span>
              </button>
            </div>
          </AwPanel>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.hello {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
  margin: 14px 0 26px;

  .eyebrow {
    font-family: var(--font-mono);
    font-size: 0.66rem;
    letter-spacing: 0.3em;
    color: var(--primary);
    text-transform: uppercase;
    margin-bottom: 10px;
    display: block;
  }
  h1 {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: clamp(1.8rem, 3.6vw, 2.7rem);
    letter-spacing: 0.02em;
    line-height: 1.05;
    text-shadow: var(--glow-text);

    em {
      font-style: normal;
      background: var(--grad-arcane);
      -webkit-background-clip: text;
      background-clip: text;
      color: transparent;
    }
  }
  p { color: var(--ink-faint); font-size: 0.9rem; margin-top: 8px; max-width: 440px; }
}

.stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 18px;
  margin-bottom: 30px;
}

.cols {
  display: grid;
  grid-template-columns: 1.7fr 1fr;
  gap: 24px;
  align-items: start;
}

.block { margin-bottom: 30px; }
.block-head {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;

  h2 { font-family: var(--font-display); font-weight: 600; font-size: 1.15rem; letter-spacing: 0.03em; }
  .rule { flex: 1; height: 1px; background: linear-gradient(90deg, var(--line-strong), transparent); }
  a { font-family: var(--font-mono); font-size: 0.68rem; color: var(--primary); text-decoration: none; letter-spacing: 0.06em; white-space: nowrap; }
}

.charscroll {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}
.card-link { display: block; text-decoration: none; color: inherit; }

.feed { padding: 6px 22px 18px; }

.quick {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 18px;

  .qa {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 15px;
    border-radius: 14px;
    border: 1px solid var(--line);
    background: var(--surface-2);
    cursor: pointer;
    transition: 0.25s;
    text-align: left;
    color: var(--ink);

    &:hover { border-color: var(--c); box-shadow: 0 0 22px -6px var(--c); transform: translateY(-3px); }
    .qi {
      width: 36px;
      height: 36px;
      border-radius: 10px;
      display: grid;
      place-items: center;
      flex: none;
      background: color-mix(in srgb, var(--c) 16%, transparent);
      color: var(--c);

      :deep(svg) { width: 18px; height: 18px; }
    }
    b { font-size: 0.84rem; display: block; }
    small { font-size: 0.66rem; color: var(--ink-faint); }
  }
}

@media (max-width: 1080px) {
  .stats { grid-template-columns: repeat(2, 1fr); }
  .cols { grid-template-columns: 1fr; }
  .charscroll { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 820px) {
  .hello { flex-direction: column; align-items: flex-start; }
  .charscroll { grid-template-columns: 1fr; }
}
</style>
