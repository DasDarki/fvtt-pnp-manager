<script setup lang="ts">
import type { ApiAsset } from '~/types/api'

const props = defineProps<{
  subjectType: 'character' | 'item' | 'scene' | 'image'
  subjectId: string
  defaultSize?: string
  promptHint?: string
}>()

const open = defineModel<boolean>('open', { default: false })
const imageUrl = defineModel<string>('imageUrl', { default: '' })

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()
const providers = useProviderStore()

const tab = ref<'library' | 'upload' | 'generate'>('library')
const assets = ref<ApiAsset[]>([])
const loading = ref(false)
const busy = ref(false)
const error = ref('')

async function load() {
  if (!campaign.currentId) return
  loading.value = true
  try {
    assets.value = await api<ApiAsset[]>(`/campaigns/${campaign.currentId}/assets`)
  } finally {
    loading.value = false
  }
}

watch(open, (v) => {
  if (v) {
    error.value = ''
    tab.value = 'library'
    prompt.value = props.promptHint || ''
    providers.ensure()
    load()
  }
})

async function attach(assetId: string, url: string) {
  if (busy.value) return
  busy.value = true
  error.value = ''
  try {
    await api(`/campaigns/${campaign.currentId}/assets/attach`, {
      method: 'POST',
      body: { subjectType: props.subjectType, subjectId: props.subjectId, assetId },
    })
    imageUrl.value = url
    open.value = false
  } catch (e: any) {
    error.value = e?.data?.error || 'Fehler'
  } finally {
    busy.value = false
  }
}

const fileInput = ref<HTMLInputElement | null>(null)
async function onFile(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file || busy.value) return
  busy.value = true
  error.value = ''
  try {
    const fd = new FormData()
    fd.append('file', file)
    const a = await api<ApiAsset>(`/campaigns/${campaign.currentId}/assets`, { method: 'POST', body: fd })
    await attach(a.id, a.url)
  } catch (e: any) {
    error.value = e?.data?.error || 'Upload fehlgeschlagen'
  } finally {
    busy.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

const prompt = ref('')
const size = ref(props.defaultSize || '1024x1024')
const sizes = [
  { value: '1024x1024', key: 'sizeSquare' },
  { value: '1792x1024', key: 'sizeWide' },
  { value: '1024x1792', key: 'sizeTall' },
]
async function generate() {
  if (!prompt.value.trim() || busy.value) return
  busy.value = true
  error.value = ''
  try {
    const r = await api<{ imageUrl: string; status: string; error: string }>(
      `/campaigns/${campaign.currentId}/dalle/generate`,
      {
        method: 'POST',
        body: {
          prompt: prompt.value,
          size: size.value,
          provider: providers.selected,
          subjectType: props.subjectType,
          subjectId: props.subjectId,
        },
      },
    )
    if (r.imageUrl) {
      imageUrl.value = r.imageUrl
      open.value = false
    } else {
      error.value = r.error || 'Fehler'
    }
  } catch (e: any) {
    error.value = e?.data?.error || 'Fehler'
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="ip-backdrop" @click="open = false" />
    <Transition name="ip">
      <div v-if="open" class="ip" role="dialog" aria-modal="true">
        <header class="ih">
          <h2><Icon name="lucide:images" /> {{ t('picker.title') }}</h2>
          <button class="x" :title="t('picker.close')" @click="open = false"><Icon name="lucide:x" /></button>
        </header>

        <div class="tabs">
          <button :class="{ on: tab === 'library' }" @click="tab = 'library'"><Icon name="lucide:library" /> {{ t('picker.library') }}</button>
          <button :class="{ on: tab === 'upload' }" @click="tab = 'upload'"><Icon name="lucide:upload" /> {{ t('picker.upload') }}</button>
          <button :class="{ on: tab === 'generate' }" @click="tab = 'generate'"><Icon name="lucide:wand-sparkles" /> {{ t('picker.generate') }}</button>
        </div>

        <div class="body">
          <template v-if="tab === 'library'">
            <div v-if="loading" class="state">{{ t('picker.loading') }}</div>
            <div v-else-if="assets.length" class="grid">
              <button v-for="a in assets" :key="a.id" class="cell" :title="a.name || a.prompt || a.source" @click="attach(a.id, a.url)">
                <img :src="a.url" :alt="a.prompt" loading="lazy" />
                <span v-if="a.source === 'upload'" class="tag up">{{ t('picker.uploaded') }}</span>
                <span v-else-if="a.source === 'mock'" class="tag mock">Mock</span>
              </button>
            </div>
            <div v-else class="state">{{ t('picker.emptyLib') }}</div>
          </template>

          <template v-else-if="tab === 'upload'">
            <button class="drop" :disabled="busy" @click="fileInput?.click()">
              <Icon name="lucide:image-up" />
              <b>{{ busy ? t('picker.uploading') : t('picker.pickFile') }}</b>
              <small>{{ t('picker.fileTypes') }}</small>
            </button>
            <input ref="fileInput" type="file" accept="image/png,image/jpeg,image/webp,image/gif,image/svg+xml" hidden @change="onFile" />
          </template>

          <template v-else>
            <textarea v-model="prompt" class="inp area" :placeholder="t('picker.promptPlaceholder')" rows="3" />
            <div class="genrow">
              <select v-model="size" class="inp">
                <option v-for="s in sizes" :key="s.value" :value="s.value">{{ t(`dalle.${s.key}`) }} · {{ s.value }}</option>
              </select>
              <ProviderPicker />
            </div>
            <AwButton icon="lucide:wand-sparkles" variant="primary" :disabled="busy || !prompt.trim()" @click="generate">
              {{ busy ? t('dalle.generating') : t('picker.generateBtn') }}
            </AwButton>
          </template>

          <p v-if="error" class="err">{{ error }}</p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.ip-backdrop { position: fixed; inset: 0; z-index: 80; background: rgba(4, 3, 10, 0.66); backdrop-filter: blur(4px); }
.ip {
  position: fixed;
  z-index: 81;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: min(720px, calc(100vw - 32px));
  max-height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 18px;
  box-shadow: 0 40px 80px -30px #000, var(--glow-secondary);
  overflow: hidden;
}

.ih {
  display: flex;
  align-items: center;
  padding: 16px 18px;
  border-bottom: 1px solid var(--line);
  h2 { display: flex; align-items: center; gap: 9px; font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; flex: 1; }
  h2 :deep(svg) { width: 18px; height: 18px; color: var(--secondary); }
  .x { display: grid; place-items: center; width: 32px; height: 32px; border-radius: 9px; border: 1px solid var(--line); background: var(--surface); color: var(--ink-faint); cursor: pointer; transition: 0.2s; }
  .x :deep(svg) { width: 16px; height: 16px; }
  .x:hover { color: var(--ember); border-color: var(--ember); }
}

.tabs {
  display: flex;
  gap: 6px;
  padding: 12px 18px 0;
  button {
    display: inline-flex;
    align-items: center;
    gap: 7px;
    padding: 9px 14px;
    border: 1px solid var(--line);
    border-bottom: 0;
    border-radius: 10px 10px 0 0;
    background: var(--surface);
    color: var(--ink-dim);
    font-size: 0.8rem;
    cursor: pointer;
    transition: 0.2s;
    :deep(svg) { width: 14px; height: 14px; }
    &:hover { color: var(--ink); }
    &.on { color: var(--primary); border-color: var(--line-strong); background: var(--surface-2); box-shadow: 0 -2px 0 var(--primary) inset; }
  }
}

.body { padding: 18px; overflow-y: auto; }

.state { padding: 40px 10px; text-align: center; color: var(--ink-faint); font-family: var(--font-mono); font-size: 0.82rem; }

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}
.cell {
  position: relative;
  border: 1px solid var(--line);
  border-radius: 12px;
  overflow: hidden;
  background: var(--void-2);
  cursor: pointer;
  padding: 0;
  transition: transform 0.2s, border-color 0.2s, box-shadow 0.2s;
  img { width: 100%; aspect-ratio: 1; object-fit: cover; display: block; }
  &:hover { transform: translateY(-3px); border-color: var(--primary); box-shadow: var(--glow-primary); }
  .tag {
    position: absolute;
    top: 6px;
    left: 6px;
    font-family: var(--font-mono);
    font-size: 0.5rem;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    padding: 2px 6px;
    border-radius: 999px;
  }
  .up { color: var(--emerald); background: rgba(55, 232, 164, 0.16); border: 1px solid var(--emerald); }
  .mock { color: var(--gold); background: rgba(255, 194, 77, 0.14); border: 1px solid var(--gold); }
}

.drop {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 44px 20px;
  border: 1.5px dashed var(--line-strong);
  border-radius: 16px;
  background: var(--surface);
  color: var(--ink-dim);
  cursor: pointer;
  transition: 0.2s;
  :deep(svg) { width: 34px; height: 34px; color: var(--secondary); }
  b { font-family: var(--font-display); font-size: 0.95rem; }
  small { font-family: var(--font-mono); font-size: 0.66rem; color: var(--ink-faint); }
  &:hover { border-color: var(--primary); color: var(--ink); }
  &:disabled { opacity: 0.6; cursor: default; }
}

.inp {
  font-family: var(--font-body);
  font-size: 0.9rem;
  color: var(--ink);
  background: var(--surface);
  border: 1px solid var(--line-strong);
  border-radius: 11px;
  padding: 11px 13px;
  width: 100%;
  transition: 0.2s;
  &:focus { outline: 0; border-color: var(--primary); box-shadow: var(--glow-primary); }
}
.area { resize: vertical; line-height: 1.5; margin-bottom: 10px; }
select.inp { cursor: pointer; appearance: none; }
.genrow { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; margin-bottom: 12px; }
:deep(.aw-btn) { width: 100%; justify-content: center; }

.err { color: var(--ember); font-size: 0.8rem; margin-top: 12px; font-family: var(--font-mono); }

.ip-enter-active, .ip-leave-active { transition: opacity 0.2s, transform 0.2s; }
.ip-enter-from, .ip-leave-to { opacity: 0; transform: translate(-50%, -48%) scale(0.98); }

@media (max-width: 560px) {
  .genrow { grid-template-columns: 1fr; }
}
</style>
