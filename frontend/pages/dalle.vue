<script setup lang="ts">
import type { ApiAsset } from '~/types/api'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()
const providers = useProviderStore()

useSeoMeta({ title: () => t('dalle.title') })

const { data: gallery, refresh } = useAsyncData(
  'dalle-assets',
  async () => {
    await campaign.ensure()
    return api<ApiAsset[]>(`/campaigns/${campaign.currentId}/assets`)
  },
  { server: false, default: () => [] },
)

onMounted(() => providers.ensure())

const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)
async function onUpload(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file || uploading.value) return
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', file)
    await api(`/campaigns/${campaign.currentId}/assets`, { method: 'POST', body: fd })
    await refresh()
  } finally {
    uploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

async function removeAsset(id: string) {
  await api(`/campaigns/${campaign.currentId}/assets/${id}`, { method: 'DELETE' })
  await refresh()
}

const stylePrompt = ref('')
const selectedPreset = ref('')
watch(
  () => campaign.current,
  (c) => {
    if (c) {
      stylePrompt.value = c.stylePrompt || ''
      selectedPreset.value = c.artStyle || ''
    }
  },
  { immediate: true },
)

function applyPreset(p: StylePreset) {
  selectedPreset.value = p.id
  stylePrompt.value = p.prompt
}

const savingStyle = ref(false)
const styleSaved = ref(false)
async function saveStyle() {
  savingStyle.value = true
  try {
    await campaign.updateStyle(selectedPreset.value, stylePrompt.value)
    styleSaved.value = true
    setTimeout(() => (styleSaved.value = false), 2000)
  } finally {
    savingStyle.value = false
  }
}

const prompt = ref('')
const size = ref('1024x1024')
const sizes = [
  { value: '1024x1024', key: 'sizeSquare' },
  { value: '1792x1024', key: 'sizeWide' },
  { value: '1024x1792', key: 'sizeTall' },
]
const generating = ref(false)
const genError = ref('')

const appliedStyle = computed(() => campaign.current?.stylePrompt?.trim() || '')

async function generate() {
  if (!prompt.value.trim() || generating.value) return
  generating.value = true
  genError.value = ''
  try {
    await api(`/campaigns/${campaign.currentId}/dalle/generate`, {
      method: 'POST',
      body: { prompt: prompt.value, size: size.value, provider: providers.selected },
    })
    await refresh()
  } catch (e: any) {
    genError.value = e?.data?.error || 'Fehler'
  } finally {
    generating.value = false
  }
}
</script>

<template>
  <div class="atelier">
    <header class="head">
      <span class="ic"><Icon name="lucide:wand-sparkles" /></span>
      <div>
        <h1>{{ t('dalle.title') }}</h1>
        <p>{{ t('dalle.subtitle') }}</p>
      </div>
    </header>

    <div class="cols">
      <div class="left">
        <AwPanel class="card">
          <div class="ch">
            <Icon name="lucide:palette" />
            <h2>{{ t('dalle.styleTitle') }}</h2>
            <span class="glow-tag">{{ t('dalle.global') }}</span>
          </div>
          <p class="hint">{{ t('dalle.styleHint') }}</p>

          <div class="presets">
            <button
              v-for="p in STYLE_PRESETS"
              :key="p.id"
              class="preset"
              :class="{ on: selectedPreset === p.id }"
              @click="applyPreset(p)"
            >
              {{ p.label }}
            </button>
          </div>

          <textarea v-model="stylePrompt" class="inp area" :placeholder="t('dalle.stylePlaceholder')" rows="3" />

          <div class="row-end">
            <AwButton icon="lucide:save" variant="soft" @click="saveStyle">
              {{ savingStyle ? t('dalle.savingStyle') : styleSaved ? t('dalle.styleSaved') : t('dalle.saveStyle') }}
            </AwButton>
          </div>
        </AwPanel>

        <AwPanel class="card">
          <div class="ch">
            <Icon name="lucide:sparkles" />
            <h2>{{ t('dalle.generateTitle') }}</h2>
            <span class="glow-tag">BYOK</span>
          </div>

          <div class="provrow">
            <span class="prl">{{ t('dalle.providerPick') }}</span>
            <select
              :value="providers.selected"
              class="inp"
              @change="providers.select(($event.target as HTMLSelectElement).value)"
            >
              <option value="">{{ t('providers.mock') }}</option>
              <option v-for="p in providers.available" :key="p.provider" :value="p.provider">{{ providerLabel(p.provider) }}</option>
            </select>
            <NuxtLink to="/settings/providers" class="manage">
              <Icon name="lucide:settings-2" /> {{ t('dalle.providerManage') }}
            </NuxtLink>
          </div>

          <textarea v-model="prompt" class="inp area" :placeholder="t('dalle.promptPlaceholder')" rows="4" />

          <div class="controls">
            <select v-model="size" class="inp">
              <option v-for="s in sizes" :key="s.value" :value="s.value">{{ t(`dalle.${s.key}`) }} · {{ s.value }}</option>
            </select>
            <AwButton icon="lucide:wand-sparkles" variant="primary" @click="generate">
              {{ generating ? t('dalle.generating') : t('dalle.generate') }}
            </AwButton>
          </div>

          <div class="composed">
            <span class="cl">{{ t('dalle.applied') }}</span>
            <p>
              <span class="p">{{ prompt || '…' }}</span>
              <span v-if="appliedStyle" class="s"> — {{ appliedStyle }}</span>
              <span v-else class="none"> · {{ t('dalle.noStyle') }}</span>
            </p>
          </div>

          <p v-if="genError" class="err">{{ genError }}</p>
        </AwPanel>
      </div>

      <div class="right">
        <div class="gh">
          <h2>{{ t('dalle.library') }}</h2>
          <span class="cnt">{{ gallery.length }}</span>
          <AwButton icon="lucide:image-up" variant="soft" :disabled="uploading" @click="fileInput?.click()">
            {{ uploading ? t('picker.uploading') : t('dalle.upload') }}
          </AwButton>
          <input
            ref="fileInput"
            type="file"
            accept="image/png,image/jpeg,image/webp,image/gif,image/svg+xml"
            hidden
            @change="onUpload"
          />
        </div>

        <div v-if="generating" class="gen-skel">
          <Icon name="lucide:loader-circle" class="spin" />
          <span>{{ t('dalle.generating') }}</span>
        </div>

        <div v-if="gallery.length" class="grid">
          <figure v-for="g in gallery" :key="g.id" class="shot">
            <span v-if="g.source === 'mock'" class="mock">{{ t('dalle.mock') }}</span>
            <span v-else-if="g.source === 'upload'" class="up">{{ t('picker.uploaded') }}</span>
            <button class="del" :title="t('actions.delete')" @click="removeAsset(g.id)"><Icon name="lucide:trash-2" /></button>
            <img :src="g.url" :alt="g.prompt" loading="lazy" />
            <figcaption>
              <p>{{ g.prompt || '—' }}</p>
              <small>{{ relativeTime(g.createdAt) }}</small>
            </figcaption>
          </figure>
        </div>
        <div v-else-if="!generating" class="empty">
          <Icon name="lucide:image" />
          {{ t('dalle.empty') }}
        </div>
      </div>
    </div>
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
    background: color-mix(in srgb, var(--secondary) 16%, transparent);
    color: var(--secondary);
    border: 1px solid var(--line);
    box-shadow: var(--glow-secondary);
    :deep(svg) { width: 24px; height: 24px; }
  }
  h1 { font-family: var(--font-display); font-weight: 600; font-size: clamp(1.6rem, 3vw, 2.2rem); letter-spacing: 0.02em; text-shadow: var(--glow-text); }
  p { color: var(--ink-faint); font-size: 0.84rem; margin-top: 4px; }
}

.cols {
  display: grid;
  grid-template-columns: 420px 1fr;
  gap: 24px;
  align-items: start;
}

.left { display: flex; flex-direction: column; gap: 18px; position: sticky; top: 90px; }
.card { padding: 22px; }

.ch {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;

  :deep(svg) { width: 17px; height: 17px; color: var(--primary); }
  h2 { font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; }
  .glow-tag {
    margin-left: auto;
    font-family: var(--font-mono);
    font-size: 0.56rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: #06040c;
    background: var(--grad-arcane);
    padding: 3px 9px;
    border-radius: 999px;
    box-shadow: var(--glow-primary);
  }
}
:global(html[data-theme='light']) .glow-tag { color: #fff; }

.hint { font-size: 0.78rem; color: var(--ink-faint); margin-bottom: 14px; }

.provrow {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  gap: 8px 12px;
  margin-bottom: 12px;

  .prl {
    grid-column: 1 / -1;
    font-family: var(--font-mono);
    font-size: 0.56rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: var(--ink-faint);
  }
  .manage {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 0.74rem;
    color: var(--ink-dim);
    text-decoration: none;
    padding: 0 6px;
    transition: 0.2s;
    :deep(svg) { width: 14px; height: 14px; }
    &:hover { color: var(--primary); }
  }
}

.presets {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 14px;

  .preset {
    font-family: var(--font-mono);
    font-size: 0.66rem;
    letter-spacing: 0.04em;
    padding: 7px 12px;
    border-radius: 9px;
    border: 1px solid var(--line);
    background: var(--surface);
    color: var(--ink-dim);
    cursor: pointer;
    transition: 0.2s;

    &:hover { color: var(--ink); border-color: var(--line-strong); }
    &.on { color: var(--secondary); border-color: var(--secondary); background: rgba(183, 104, 255, 0.1); box-shadow: var(--glow-secondary); }
  }
}

.inp {
  font-family: var(--font-body);
  font-size: 0.92rem;
  color: var(--ink);
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 12px;
  padding: 12px 14px;
  width: 100%;
  transition: 0.25s;

  &:focus { outline: 0; border-color: var(--primary); box-shadow: 0 0 0 3px rgba(70, 232, 255, 0.14), var(--glow-primary); }
}
.area { resize: vertical; line-height: 1.5; }
select.inp { cursor: pointer; appearance: none; }

.row-end { display: flex; justify-content: flex-end; margin-top: 14px; }

.controls {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 10px;
  margin-top: 12px;
}

.composed {
  margin-top: 14px;
  padding: 13px 15px;
  border-radius: 12px;
  border: 1px dashed var(--line-strong);
  background: var(--surface);

  .cl { font-family: var(--font-mono); font-size: 0.56rem; letter-spacing: 0.16em; text-transform: uppercase; color: var(--ink-faint); }
  p { font-size: 0.82rem; margin-top: 6px; line-height: 1.5; }
  .p { color: var(--ink); }
  .s { color: var(--secondary); }
  .none { color: var(--ink-faint); font-style: italic; }
}

.err { color: var(--ember); font-size: 0.8rem; margin-top: 12px; font-family: var(--font-mono); }

.gh {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;

  h2 { font-family: var(--font-display); font-weight: 600; font-size: 1.15rem; }
  .cnt { font-family: var(--font-mono); font-size: 0.7rem; color: var(--ink-faint); border: 1px solid var(--line); border-radius: 999px; padding: 3px 10px; }
  :deep(.aw-btn) { margin-left: auto; }
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
}

.shot {
  position: relative;
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--line);
  background: var(--surface);
  box-shadow: var(--shadow-panel);
  transition: transform 0.3s, box-shadow 0.3s;

  &:hover { transform: translateY(-4px); box-shadow: 0 26px 50px -26px #000, var(--glow-secondary); }
  &:hover .del { opacity: 1; }
  img { width: 100%; aspect-ratio: 1; object-fit: cover; display: block; background: var(--void-2); }
  .mock,
  .up {
    position: absolute;
    top: 10px;
    left: 10px;
    z-index: 2;
    font-family: var(--font-mono);
    font-size: 0.52rem;
    letter-spacing: 0.12em;
    padding: 3px 8px;
    border-radius: 999px;
  }
  .mock { color: var(--gold); background: rgba(255, 194, 77, 0.14); border: 1px solid var(--gold); }
  .up { color: var(--emerald); background: rgba(55, 232, 164, 0.16); border: 1px solid var(--emerald); }
  .del {
    position: absolute;
    top: 8px;
    right: 8px;
    z-index: 2;
    display: grid;
    place-items: center;
    width: 28px;
    height: 28px;
    border-radius: 8px;
    border: 1px solid var(--line-strong);
    background: rgba(8, 6, 14, 0.7);
    color: var(--ink-dim);
    cursor: pointer;
    opacity: 0;
    transition: 0.2s;
    :deep(svg) { width: 14px; height: 14px; }
    &:hover { color: var(--ember); border-color: var(--ember); }
  }
  figcaption {
    padding: 12px 14px;
    p { font-size: 0.78rem; color: var(--ink-dim); display: -webkit-box; -webkit-line-clamp: 2; line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
    small { font-family: var(--font-mono); font-size: 0.62rem; color: var(--ink-faint); display: block; margin-top: 6px; }
  }
}

.gen-skel,
.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 50px 20px;
  color: var(--ink-faint);
  font-family: var(--font-mono);
  font-size: 0.8rem;
  border: 1.5px dashed var(--line-strong);
  border-radius: 18px;
  margin-bottom: 16px;
  :deep(svg) { width: 30px; height: 30px; opacity: 0.7; }
}
.spin { animation: spin 0.9s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 1000px) {
  .cols { grid-template-columns: 1fr; }
  .left { position: static; }
}
</style>
