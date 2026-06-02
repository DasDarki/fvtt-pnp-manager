<script setup lang="ts">
withDefaults(defineProps<{ empty?: boolean; min?: number }>(), { empty: false, min: 260 })
</script>

<template>
  <div v-if="empty" class="empty">
    <Icon name="lucide:search-x" />
    <slot name="empty">Nichts gefunden.</slot>
  </div>
  <div v-else class="grid" :style="{ '--min': `${min}px` }">
    <slot />
  </div>
</template>

<style lang="scss" scoped>
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(var(--min), 1fr));
  gap: 18px;
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding: 70px 20px;
  text-align: center;
  color: var(--ink-faint);
  font-family: var(--font-mono);
  font-size: 0.82rem;
  letter-spacing: 0.06em;
  border: 1.5px dashed var(--line-strong);
  border-radius: 22px;

  :deep(svg) { width: 34px; height: 34px; color: var(--ink-faint); opacity: 0.7; }
}
</style>
