<script setup lang="ts">
interface Folder {
  id: string
  name: string
  foundryType: string
  parentId: string | null
}

const props = defineProps<{ type: string }>()
const model = defineModel<string | null>()

const { t } = useI18n()
const api = useApi()
const campaign = useCampaignStore()
const folders = ref<Folder[]>([])

onMounted(async () => {
  await campaign.ensure()
  try {
    folders.value = await api<Folder[]>(`/campaigns/${campaign.currentId}/folders`)
  } catch {
    folders.value = []
  }
})

const options = computed(() => {
  const list = folders.value.filter((f) => f.foundryType === props.type)
  const ids = new Set(list.map((f) => f.id))
  const children: Record<string, Folder[]> = {}
  for (const f of list) {
    const key = f.parentId && ids.has(f.parentId) ? f.parentId : '__root'
    ;(children[key] ??= []).push(f)
  }
  const rows: { id: string; label: string }[] = []
  const walk = (pid: string, depth: number) => {
    for (const f of (children[pid] || []).slice().sort((a, b) => a.name.localeCompare(b.name))) {
      rows.push({ id: f.id, label: `${'  '.repeat(depth)}${depth ? '↳ ' : ''}${f.name}` })
      walk(f.id, depth + 1)
    }
  }
  walk('__root', 0)
  return rows
})

function onChange(e: Event) {
  const v = (e.target as HTMLSelectElement).value
  model.value = v || null
}
</script>

<template>
  <select class="fp" :value="model || ''" :disabled="!options.length" @change="onChange">
    <option value="">{{ options.length ? t('folder.none') : t('folder.empty') }}</option>
    <option v-for="o in options" :key="o.id" :value="o.id">{{ o.label }}</option>
  </select>
</template>

<style lang="scss" scoped>
.fp {
  font-family: var(--font-body);
  font-size: 0.92rem;
  color: var(--ink);
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 11px;
  padding: 11px 13px;
  width: 100%;
  cursor: pointer;
  appearance: none;
  transition: 0.25s;

  &:disabled { opacity: 0.55; cursor: not-allowed; }
  &:focus { outline: 0; border-color: var(--primary); box-shadow: 0 0 0 3px rgba(70, 232, 255, 0.14), var(--glow-primary); }
}
</style>
