<script setup lang="ts">
const { t } = useI18n()
const campaign = useCampaignStore()

const open = ref(false)
const creating = ref(false)
const showCreate = ref(false)
const newName = ref('')
const newRuleset = ref('dnd5e_2024')

const rulesets = [
  { value: 'dnd5e_2024', label: 'D&D 5e (2024)' },
  { value: 'vampire_v5', label: 'Vampire V5' },
]
function rulesetLabel(r?: string) {
  return rulesets.find((x) => x.value === r)?.label || '—'
}

function switchTo(id: string) {
  if (id === campaign.currentId) {
    open.value = false
    return
  }
  campaign.select(id)
  window.location.href = '/dashboard'
}

async function create() {
  if (!newName.value.trim() || creating.value) return
  creating.value = true
  try {
    const cam = await campaign.createCampaign(newName.value.trim(), newRuleset.value)
    campaign.select(cam.id)
    window.location.href = '/dashboard'
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <div class="cs">
    <button class="campaign" @click="open = !open">
      <span class="cv" aria-hidden="true" />
      <span class="ct">
        <b>{{ campaign.current?.name || 'Aetherwright' }}</b>
        <small>{{ rulesetLabel(campaign.current?.ruleset) }}</small>
      </span>
      <Icon name="lucide:chevrons-up-down" class="chev" />
    </button>

    <Teleport to="body">
      <div v-if="open" class="cs-backdrop" @click="open = false" />
    </Teleport>

    <Transition name="drop">
      <div v-if="open" class="menu">
        <div class="grp">{{ t('campaign.switch') }}</div>
        <button
          v-for="c in campaign.campaigns"
          :key="c.id"
          class="row"
          :class="{ on: c.id === campaign.currentId }"
          @click="switchTo(c.id)"
        >
          <span class="dot" />
          <span class="rt">
            <b>{{ c.name }}</b>
            <small>{{ rulesetLabel(c.ruleset) }}</small>
          </span>
          <Icon v-if="c.id === campaign.currentId" name="lucide:check" class="ck" />
        </button>

        <div v-if="!showCreate" class="foot">
          <button class="act" @click="showCreate = true"><Icon name="lucide:plus" /> {{ t('campaign.new') }}</button>
          <NuxtLink to="/settings/campaigns" class="act" @click="open = false"><Icon name="lucide:settings-2" /> {{ t('campaign.manage') }}</NuxtLink>
        </div>

        <div v-else class="create">
          <input v-model="newName" class="inp" :placeholder="t('campaign.namePlaceholder')" @keydown.enter="create" />
          <select v-model="newRuleset" class="inp">
            <option v-for="r in rulesets" :key="r.value" :value="r.value">{{ r.label }}</option>
          </select>
          <AwButton icon="lucide:sparkles" variant="primary" @click="create">
            {{ creating ? t('campaign.creating') : t('campaign.create') }}
          </AwButton>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style lang="scss" scoped>
.cs { position: relative; margin: 0 16px 6px; }

.campaign {
  width: 100%;
  padding: 13px 14px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: var(--surface-2);
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  text-align: left;
  transition: 0.25s;
  color: var(--ink);

  &:hover { border-color: var(--line-strong); }
  .cv {
    width: 36px;
    height: 36px;
    border-radius: 10px;
    flex: none;
    background: conic-gradient(from 120deg, var(--secondary), var(--magenta), var(--gold), var(--secondary));
    box-shadow: var(--glow-secondary);
  }
  .ct { flex: 1; min-width: 0; }
  b { font-family: var(--font-display); font-weight: 600; font-size: 0.92rem; display: block; line-height: 1.1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  small { font-size: 0.68rem; color: var(--ink-faint); }
  .chev { color: var(--ink-faint); width: 16px; height: 16px; }
}

.cs-backdrop { position: fixed; inset: 0; z-index: 54; }

.menu {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  z-index: 55;
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 14px;
  box-shadow: 0 30px 60px -24px #000, var(--glow-secondary);
  backdrop-filter: blur(14px);
  padding: 8px;
  overflow: hidden;
}

.grp {
  font-family: var(--font-mono);
  font-size: 0.56rem;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: var(--ink-faint);
  padding: 8px 10px 6px;
}

.row {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 9px 10px;
  border-radius: 10px;
  border: 0;
  background: transparent;
  cursor: pointer;
  text-align: left;
  transition: 0.18s;

  &:hover { background: var(--surface); }
  &.on { background: rgba(70, 232, 255, 0.08); }
  .dot { width: 8px; height: 8px; border-radius: 50%; flex: none; background: var(--secondary); box-shadow: 0 0 8px var(--secondary); }
  .rt { flex: 1; min-width: 0; }
  b { font-size: 0.86rem; font-weight: 600; color: var(--ink); display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  small { font-size: 0.66rem; color: var(--ink-faint); }
  .ck { width: 15px; height: 15px; color: var(--primary); }
}

.foot {
  display: flex;
  gap: 6px;
  margin-top: 6px;
  padding-top: 8px;
  border-top: 1px solid var(--line);

  .act {
    flex: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    padding: 9px;
    border-radius: 9px;
    border: 1px solid var(--line);
    background: var(--surface);
    color: var(--ink-dim);
    font-size: 0.74rem;
    cursor: pointer;
    text-decoration: none;
    transition: 0.2s;
    :deep(svg) { width: 14px; height: 14px; }
    &:hover { color: var(--primary); border-color: var(--primary); }
  }
}

.create {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
  padding-top: 10px;
  border-top: 1px solid var(--line);

  .inp {
    font-family: var(--font-body);
    font-size: 0.86rem;
    color: var(--ink);
    background: var(--surface);
    border: 1px solid var(--line-strong);
    border-radius: 9px;
    padding: 9px 11px;
    transition: 0.2s;
    &:focus { outline: 0; border-color: var(--primary); box-shadow: var(--glow-primary); }
  }
  select.inp { cursor: pointer; appearance: none; }
  :deep(.aw-btn) { width: 100%; justify-content: center; }
}

.drop-enter-active, .drop-leave-active { transition: opacity 0.18s, transform 0.18s; }
.drop-enter-from, .drop-leave-to { opacity: 0; transform: translateY(-6px); }
</style>
