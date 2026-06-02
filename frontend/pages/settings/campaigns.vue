<script setup lang="ts">
import type { ApiCampaign } from '~/types/api'

definePageMeta({ layout: 'deck' })

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()

useSeoMeta({ title: () => t('campaign.title') })

const rulesets = [
  { value: 'dnd5e_2024', label: 'D&D 5e (2024)' },
  { value: 'vampire_v5', label: 'Vampire V5' },
]

const newName = ref('')
const newRuleset = ref('dnd5e_2024')
const creating = ref(false)

onMounted(() => campaign.ensure())

async function create() {
  if (!newName.value.trim() || creating.value) return
  creating.value = true
  try {
    const cam = await campaign.createCampaign(newName.value.trim(), newRuleset.value)
    newName.value = ''
    campaign.select(cam.id)
    window.location.href = '/dashboard'
  } finally {
    creating.value = false
  }
}

const savedId = ref('')
async function saveRow(c: ApiCampaign) {
  await api(`/campaigns/${c.id}`, {
    method: 'PATCH',
    body: { name: c.name, ruleset: c.ruleset, forgeRootPath: c.forgeRootPath },
  })
  await campaign.refreshList()
  savedId.value = c.id
  setTimeout(() => (savedId.value = ''), 2000)
}

function setActive(id: string) {
  if (id === campaign.currentId) return
  campaign.select(id)
  window.location.href = '/dashboard'
}

const confirmingId = ref('')
async function del(c: ApiCampaign) {
  if (confirmingId.value !== c.id) {
    confirmingId.value = c.id
    setTimeout(() => (confirmingId.value = ''), 3000)
    return
  }
  const wasCurrent = c.id === campaign.currentId
  await api(`/campaigns/${c.id}`, { method: 'DELETE' })
  await campaign.refreshList()
  if (wasCurrent) window.location.href = '/dashboard'
}
</script>

<template>
  <div class="cm">
    <header class="head">
      <span class="ic"><Icon name="lucide:library" /></span>
      <div>
        <h1>{{ t('campaign.title') }}</h1>
        <p>{{ t('campaign.subtitle') }}</p>
      </div>
    </header>

    <AwPanel class="card create">
      <div class="ch"><Icon name="lucide:plus" /><h2>{{ t('campaign.new') }}</h2></div>
      <div class="crow">
        <input v-model="newName" class="inp" :placeholder="t('campaign.namePlaceholder')" @keydown.enter="create" />
        <select v-model="newRuleset" class="inp rs">
          <option v-for="r in rulesets" :key="r.value" :value="r.value">{{ r.label }}</option>
        </select>
        <AwButton icon="lucide:sparkles" variant="primary" @click="create">
          {{ creating ? t('campaign.creating') : t('campaign.create') }}
        </AwButton>
      </div>
    </AwPanel>

    <div class="list">
      <AwPanel v-for="c in campaign.campaigns" :key="c.id" class="card row" :class="{ active: c.id === campaign.currentId }">
        <div class="grid">
          <div class="f">
            <label>{{ t('editor.name') }}</label>
            <input v-model="c.name" class="inp" type="text" />
          </div>
          <div class="f">
            <label>{{ t('campaign.ruleset') }}</label>
            <select v-model="c.ruleset" class="inp">
              <option v-for="r in rulesets" :key="r.value" :value="r.value">{{ r.label }}</option>
            </select>
          </div>
          <div class="f wide">
            <label>{{ t('campaign.forgeRoot') }}</label>
            <input v-model="c.forgeRootPath" class="inp" type="text" placeholder="worlds/…" />
          </div>
        </div>
        <div class="acts">
          <span v-if="c.id === campaign.currentId" class="badge"><Icon name="lucide:check" /> {{ t('campaign.active') }}</span>
          <AwButton v-else icon="lucide:log-in" variant="ghost" @click="setActive(c.id)">{{ t('campaign.setActive') }}</AwButton>
          <AwButton icon="lucide:save" variant="soft" @click="saveRow(c)">{{ savedId === c.id ? t('editor.saved') : t('editor.save') }}</AwButton>
          <AwButton :icon="confirmingId === c.id ? 'lucide:trash-2' : 'lucide:trash'" variant="danger" @click="del(c)">
            {{ confirmingId === c.id ? t('editor.confirmDelete') : t('editor.delete') }}
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

.card { padding: 20px; }
.ch { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; :deep(svg) { width: 17px; height: 17px; color: var(--primary); } h2 { font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; } }

.create { margin-bottom: 22px; }
.crow { display: grid; grid-template-columns: 1fr 220px auto; gap: 12px; }

.list { display: flex; flex-direction: column; gap: 14px; }
.row.active { border-color: color-mix(in srgb, var(--primary) 45%, transparent); box-shadow: 0 0 26px -10px var(--primary), var(--shadow-panel); }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px; }
.f { display: flex; flex-direction: column; gap: 7px; }
.f.wide { grid-column: 1 / -1; }
.f label { font-family: var(--font-mono); font-size: 0.6rem; letter-spacing: 0.14em; text-transform: uppercase; color: var(--ink-faint); }

.inp {
  font-family: var(--font-body);
  font-size: 0.9rem;
  color: var(--ink);
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 10px;
  padding: 10px 12px;
  transition: 0.25s;
  &:focus { outline: 0; border-color: var(--primary); box-shadow: 0 0 0 3px rgba(70, 232, 255, 0.14), var(--glow-primary); }
}
select.inp { cursor: pointer; appearance: none; }

.acts { display: flex; align-items: center; gap: 9px; flex-wrap: wrap; }
.badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 0.66rem;
  letter-spacing: 0.06em;
  color: var(--primary);
  border: 1px solid var(--primary);
  border-radius: 999px;
  padding: 7px 13px;
  :deep(svg) { width: 13px; height: 13px; }
}

@media (max-width: 720px) {
  .crow { grid-template-columns: 1fr; }
  .grid { grid-template-columns: 1fr; }
}
</style>
