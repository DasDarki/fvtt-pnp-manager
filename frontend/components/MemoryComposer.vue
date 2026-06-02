<script setup lang="ts">
import type { MemoryLevel, SubjectType } from '~/types/entities'

const props = defineProps<{
  subjectType?: SubjectType
  subjectId?: string
  subjectLabel?: string
  onSaved?: () => void | Promise<void>
}>()

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()

const title = ref('')
const body = ref('')
const level = ref<MemoryLevel>('info')
const busy = ref(false)

async function save() {
  if (!title.value.trim() || busy.value || !campaign.currentId) return
  busy.value = true
  try {
    await api(`/campaigns/${campaign.currentId}/memories`, {
      method: 'POST',
      body: {
        title: title.value.trim(),
        body: body.value.trim(),
        level: level.value,
        subjectType: props.subjectType || 'campaign',
        subjectId: props.subjectId || null,
        subjectLabel: props.subjectLabel || '',
      },
    })
    title.value = ''
    body.value = ''
    level.value = 'info'
    await props.onSaved?.()
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <AwPanel class="composer">
    <input v-model="title" class="inp" :placeholder="t('memory.titlePlaceholder')" @keydown.enter="save" />
    <textarea v-model="body" class="inp area" rows="2" :placeholder="t('memory.bodyPlaceholder')" />
    <div class="row">
      <LevelPicker v-model="level" />
      <AwButton icon="lucide:plus" variant="primary" @click="save">
        {{ busy ? t('memory.saving') : t('memory.add') }}
      </AwButton>
    </div>
  </AwPanel>
</template>

<style lang="scss" scoped>
.composer {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 11px;
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
  width: 100%;
  &:focus { outline: 0; border-color: var(--primary); box-shadow: 0 0 0 3px rgba(70, 232, 255, 0.14), var(--glow-primary); }
}
.area { resize: vertical; line-height: 1.5; }
.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}
</style>
