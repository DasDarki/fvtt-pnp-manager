<script setup lang="ts">
const { t } = useI18n()
const store = useProviderStore()

onMounted(() => store.ensure())

function onChange(e: Event) {
  store.select((e.target as HTMLSelectElement).value)
}
</script>

<template>
  <div class="pp">
    <Icon name="lucide:cpu" class="pi" />
    <select :value="store.selected" class="ps" :title="t('dalle.providerPick')" @change="onChange">
      <option value="">{{ t('providers.mock') }}</option>
      <option v-for="p in store.available" :key="p.provider" :value="p.provider">{{ providerLabel(p.provider) }}</option>
    </select>
    <NuxtLink to="/settings/providers" class="addk" :title="t('providers.title')">
      <Icon name="lucide:settings-2" />
    </NuxtLink>
  </div>
</template>

<style lang="scss" scoped>
.pp {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 4px 4px 10px;
  border: 1px solid var(--line-strong);
  border-radius: 11px;
  background: var(--surface-2);

  .pi { width: 15px; height: 15px; color: var(--primary); flex: none; }
  .ps {
    flex: 1;
    min-width: 0;
    appearance: none;
    border: 0;
    background: transparent;
    color: var(--ink);
    font-family: var(--font-body);
    font-size: 0.82rem;
    padding: 6px 2px;
    cursor: pointer;
    &:focus { outline: 0; }
  }
  .addk {
    flex: none;
    display: grid;
    place-items: center;
    width: 30px;
    height: 30px;
    border-radius: 8px;
    border: 1px solid var(--line);
    background: var(--surface);
    color: var(--ink-faint);
    transition: 0.2s;
    :deep(svg) { width: 15px; height: 15px; }
    &:hover { color: var(--primary); border-color: var(--primary); }
  }
}
</style>
