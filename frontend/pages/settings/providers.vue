<script setup lang="ts">
definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const store = useProviderStore()

useSeoMeta({ title: () => t('providers.title') })

interface Form {
  model: string
  apiKey: string
  saving: boolean
  saved: boolean
}
const forms = reactive<Record<string, Form>>({})

function init() {
  for (const p of store.catalog) {
    const s = store.providers.find((x) => x.provider === p)
    forms[p] = {
      model: s?.model || PROVIDER_META[p]?.defaultModel || '',
      apiKey: '',
      saving: false,
      saved: false,
    }
  }
}

onMounted(async () => {
  await store.ensure()
  init()
})

function hasKey(p: string) {
  return !!store.providers.find((x) => x.provider === p)?.hasKey
}

async function save(p: string) {
  const f = forms[p]
  if (f.saving) return
  f.saving = true
  try {
    await store.save(p, f.model, f.apiKey)
    f.apiKey = ''
    f.saved = true
    setTimeout(() => (f.saved = false), 2000)
  } finally {
    f.saving = false
  }
}

async function remove(p: string) {
  await store.remove(p)
  forms[p].model = PROVIDER_META[p]?.defaultModel || ''
  forms[p].apiKey = ''
}

const defaultLabel = computed(() => (store.selected ? providerLabel(store.selected) : t('providers.mock')))
</script>

<template>
  <div class="providers">
    <header class="head">
      <span class="ic"><Icon name="lucide:key-round" /></span>
      <div>
        <h1>{{ t('providers.title') }}</h1>
        <p>{{ t('providers.subtitle') }}</p>
      </div>
    </header>

    <AwPanel class="defcard">
      <div class="dl">
        <Icon name="lucide:star" />
        <div>
          <b>{{ t('providers.defaultTitle') }}</b>
          <small>{{ t('providers.defaultHint') }}</small>
        </div>
      </div>
      <select :value="store.selected" class="inp def" @change="store.select(($event.target as HTMLSelectElement).value)">
        <option value="">{{ t('providers.mock') }}</option>
        <option v-for="p in store.available" :key="p.provider" :value="p.provider">{{ providerLabel(p.provider) }}</option>
      </select>
      <span class="now">{{ defaultLabel }}</span>
    </AwPanel>

    <div class="grid">
      <AwPanel v-for="p in store.catalog" :key="p" class="card" :class="{ active: hasKey(p) }">
        <div class="ch">
          <span class="pic"><Icon :name="providerIcon(p)" /></span>
          <h2>{{ providerLabel(p) }}</h2>
          <span class="badge" :class="{ on: hasKey(p) }">
            {{ hasKey(p) ? t('providers.configured') : t('providers.notConfigured') }}
          </span>
        </div>

        <p class="host" v-if="PROVIDER_META[p]"><Icon name="lucide:external-link" /> {{ PROVIDER_META[p].keyHost }}</p>

        <label class="fl">
          <span>{{ t('providers.model') }}</span>
          <input v-if="forms[p]" v-model="forms[p].model" class="inp" :placeholder="PROVIDER_META[p]?.defaultModel" />
        </label>

        <label class="fl">
          <span>{{ t('providers.key') }}</span>
          <input
            v-if="forms[p]"
            v-model="forms[p].apiKey"
            class="inp"
            type="password"
            autocomplete="off"
            :placeholder="hasKey(p) ? t('providers.keySet') : t('providers.keyPlaceholder')"
          />
        </label>

        <p class="note">{{ t('providers.byokNote') }}</p>

        <div class="actions">
          <AwButton
            v-if="hasKey(p) && store.selected !== p"
            icon="lucide:star"
            variant="ghost"
            @click="store.select(p)"
          >
            {{ t('providers.makeDefault') }}
          </AwButton>
          <span v-else-if="store.selected === p" class="isdef"><Icon name="lucide:star" /> {{ t('providers.default') }}</span>
          <span class="sp" />
          <AwButton v-if="hasKey(p)" icon="lucide:trash-2" variant="ghost" @click="remove(p)">
            {{ t('providers.remove') }}
          </AwButton>
          <AwButton v-if="forms[p]" icon="lucide:save" variant="soft" :disabled="forms[p].saving" @click="save(p)">
            {{ forms[p].saving ? t('providers.saving') : forms[p].saved ? t('providers.saved') : t('providers.save') }}
          </AwButton>
        </div>
      </AwPanel>
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

.defcard {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  margin-bottom: 22px;

  .dl {
    display: flex;
    align-items: center;
    gap: 12px;
    :deep(svg) { width: 18px; height: 18px; color: var(--gold); }
    b { font-family: var(--font-display); font-weight: 600; font-size: 0.96rem; display: block; }
    small { font-size: 0.74rem; color: var(--ink-faint); }
  }
  .def { width: auto; min-width: 230px; margin-left: auto; }
  .now {
    font-family: var(--font-mono);
    font-size: 0.7rem;
    color: var(--secondary);
    border: 1px solid var(--line);
    border-radius: 999px;
    padding: 5px 12px;
  }
}

.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 20px; align-items: start; }
.card {
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 13px;
  transition: 0.25s;
  &.active { border-color: rgba(55, 232, 164, 0.35); }
}

.ch {
  display: flex;
  align-items: center;
  gap: 11px;

  .pic {
    width: 34px;
    height: 34px;
    border-radius: 10px;
    flex: none;
    display: grid;
    place-items: center;
    background: var(--surface-2);
    border: 1px solid var(--line);
    :deep(svg) { width: 17px; height: 17px; color: var(--primary); }
  }
  h2 { font-family: var(--font-display); font-weight: 600; font-size: 1rem; flex: 1; min-width: 0; }
  .badge {
    font-family: var(--font-mono);
    font-size: 0.56rem;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    padding: 4px 9px;
    border-radius: 999px;
    color: var(--ink-faint);
    border: 1px solid var(--line);
    &.on { color: var(--emerald); border-color: rgba(55, 232, 164, 0.4); background: rgba(55, 232, 164, 0.08); }
  }
}

.host {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 0.7rem;
  color: var(--ink-faint);
  margin-top: -4px;
  :deep(svg) { width: 12px; height: 12px; }
}

.fl {
  display: flex;
  flex-direction: column;
  gap: 6px;
  > span { font-family: var(--font-mono); font-size: 0.58rem; letter-spacing: 0.12em; text-transform: uppercase; color: var(--ink-faint); }
}

.inp {
  font-family: var(--font-body);
  font-size: 0.9rem;
  color: var(--ink);
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 11px;
  padding: 11px 13px;
  width: 100%;
  transition: 0.25s;
  &:focus { outline: 0; border-color: var(--primary); box-shadow: 0 0 0 3px rgba(70, 232, 255, 0.14), var(--glow-primary); }
}
select.inp { cursor: pointer; appearance: none; }

.note { font-size: 0.66rem; color: var(--ink-faint); font-family: var(--font-mono); line-height: 1.5; }

.actions {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 4px;
  .sp { flex: 1; }
  .isdef {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-family: var(--font-mono);
    font-size: 0.68rem;
    color: var(--gold);
    :deep(svg) { width: 13px; height: 13px; }
  }
}

@media (max-width: 720px) {
  .defcard { flex-wrap: wrap; .def { margin-left: 0; width: 100%; } }
}
</style>
