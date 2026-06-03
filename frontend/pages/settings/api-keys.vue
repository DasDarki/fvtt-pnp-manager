<script setup lang="ts">
definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const api = useApi()

useSeoMeta({ title: () => t('apikeys.title') })

interface ApiKey {
  id: string
  name: string
  prefix: string
  scope: string
  lastUsedAt: string | null
  createdAt: string
}

const { data: keys, refresh } = useAsyncData('apikeys', () => api<ApiKey[]>('/apikeys'), {
  server: false,
  default: () => [] as ApiKey[],
})

const name = ref('')
const creating = ref(false)
const freshKey = ref('')
const copied = ref(false)

async function create() {
  if (creating.value) return
  creating.value = true
  freshKey.value = ''
  try {
    const r = await api<{ key: string }>('/apikeys', { method: 'POST', body: { name: name.value.trim() } })
    freshKey.value = r.key
    name.value = ''
    await refresh()
  } finally {
    creating.value = false
  }
}

async function copyKey() {
  try {
    await navigator.clipboard.writeText(freshKey.value)
    copied.value = true
    setTimeout(() => (copied.value = false), 1500)
  } catch {
    /* ignore */
  }
}

async function revoke(id: string) {
  await api(`/apikeys/${id}`, { method: 'DELETE' })
  await refresh()
}
</script>

<template>
  <div class="keys">
    <header class="head">
      <span class="ic"><Icon name="lucide:key-square" /></span>
      <div>
        <h1>{{ t('apikeys.title') }}</h1>
        <p>{{ t('apikeys.subtitle') }}</p>
      </div>
    </header>

    <AwPanel class="card create">
      <div class="ch"><Icon name="lucide:plus" /><h2>{{ t('apikeys.newTitle') }}</h2><span class="scope">images</span></div>
      <p class="hint">{{ t('apikeys.newHint') }}</p>
      <div class="row">
        <input v-model="name" class="inp" :placeholder="t('apikeys.namePlaceholder')" @keydown.enter="create" />
        <AwButton icon="lucide:sparkles" variant="primary" :disabled="creating" @click="create">
          {{ creating ? t('apikeys.creating') : t('apikeys.create') }}
        </AwButton>
      </div>

      <div v-if="freshKey" class="fresh">
        <p class="once">{{ t('apikeys.once') }}</p>
        <button class="keyfield" @click="copyKey">
          <code>{{ freshKey }}</code>
          <Icon :name="copied ? 'lucide:check' : 'lucide:copy'" />
        </button>
      </div>
    </AwPanel>

    <AwPanel class="card">
      <div class="ch"><Icon name="lucide:list" /><h2>{{ t('apikeys.listTitle') }}</h2><span class="cnt">{{ keys.length }}</span></div>
      <div v-if="keys.length" class="list">
        <div v-for="k in keys" :key="k.id" class="krow">
          <span class="kk"><Icon name="lucide:key-round" /></span>
          <div class="km">
            <b>{{ k.name }}</b>
            <small><code>{{ k.prefix }}…</code> · {{ t('apikeys.scope') }}: {{ k.scope }}</small>
          </div>
          <span class="used">{{ k.lastUsedAt ? t('apikeys.usedAt', { t: relativeTime(k.lastUsedAt) }) : t('apikeys.neverUsed') }}</span>
          <button class="revoke" :title="t('apikeys.revoke')" @click="revoke(k.id)"><Icon name="lucide:trash-2" /></button>
        </div>
      </div>
      <p v-else class="empty">{{ t('apikeys.empty') }}</p>
    </AwPanel>

    <AwPanel class="card guide">
      <div class="ch"><Icon name="lucide:puzzle" /><h2>{{ t('apikeys.extTitle') }}</h2></div>
      <ol class="steps">
        <li>{{ t('apikeys.ext1') }}</li>
        <li>{{ t('apikeys.ext2') }}</li>
        <li>{{ t('apikeys.ext3') }}</li>
      </ol>
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
    background: color-mix(in srgb, var(--primary) 16%, transparent);
    color: var(--primary);
    border: 1px solid var(--line);
    box-shadow: var(--glow-primary);
    :deep(svg) { width: 24px; height: 24px; }
  }
  h1 { font-family: var(--font-display); font-weight: 600; font-size: clamp(1.6rem, 3vw, 2.2rem); letter-spacing: 0.02em; text-shadow: var(--glow-text); }
  p { color: var(--ink-faint); font-size: 0.84rem; margin-top: 4px; }
}

.card { padding: 22px; display: flex; flex-direction: column; gap: 14px; margin-bottom: 20px; }
.ch {
  display: flex;
  align-items: center;
  gap: 10px;
  :deep(svg) { width: 17px; height: 17px; color: var(--primary); }
  h2 { font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; }
  .scope, .cnt { margin-left: auto; font-family: var(--font-mono); font-size: 0.62rem; color: var(--ink-faint); border: 1px solid var(--line); border-radius: 999px; padding: 3px 10px; }
}
.hint { font-size: 0.8rem; color: var(--ink-faint); }

.row { display: grid; grid-template-columns: 1fr auto; gap: 10px; }
.inp {
  font-family: var(--font-body);
  font-size: 0.9rem;
  color: var(--ink);
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 11px;
  padding: 11px 13px;
  &:focus { outline: 0; border-color: var(--primary); box-shadow: var(--glow-primary); }
}

.fresh {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px;
  border: 1px dashed var(--line-strong);
  border-radius: 12px;
  background: var(--surface);
  .once { font-family: var(--font-mono); font-size: 0.66rem; color: var(--gold); }
  .keyfield {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 13px;
    border-radius: 10px;
    border: 1px solid var(--line-strong);
    background: var(--surface-2);
    cursor: pointer;
    text-align: left;
    &:hover { border-color: var(--primary); }
    code { flex: 1; font-family: var(--font-mono); font-size: 0.76rem; color: var(--ink); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
    :deep(svg) { width: 15px; height: 15px; color: var(--primary); flex: none; }
  }
}

.list { display: flex; flex-direction: column; gap: 8px; }
.krow {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 11px 13px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: var(--surface-2);
  .kk { display: grid; place-items: center; width: 32px; height: 32px; border-radius: 9px; background: var(--surface); border: 1px solid var(--line); flex: none; :deep(svg) { width: 15px; height: 15px; color: var(--secondary); } }
  .km { flex: 1; min-width: 0; b { font-size: 0.88rem; display: block; } small { font-size: 0.7rem; color: var(--ink-faint); } small code { font-family: var(--font-mono); } }
  .used { font-family: var(--font-mono); font-size: 0.66rem; color: var(--ink-faint); }
  .revoke { flex: none; display: grid; place-items: center; width: 32px; height: 32px; border-radius: 9px; border: 1px solid var(--line); background: var(--surface); color: var(--ink-faint); cursor: pointer; transition: 0.2s; :deep(svg) { width: 15px; height: 15px; } &:hover { color: var(--ember); border-color: var(--ember); } }
}
.empty { font-size: 0.82rem; color: var(--ink-faint); font-family: var(--font-mono); }

.guide .steps {
  margin: 0;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 0.84rem;
  color: var(--ink-dim);
  li::marker { color: var(--primary); font-family: var(--font-mono); }
}

@media (max-width: 640px) {
  .krow { flex-wrap: wrap; .used { order: 4; width: 100%; } }
}
</style>
