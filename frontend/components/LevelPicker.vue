<script setup lang="ts">
import type { MemoryLevel } from '~/types/entities'

const model = defineModel<MemoryLevel>({ default: 'info' })
const { t } = useI18n()
const levels: MemoryLevel[] = ['info', 'notice', 'warning', 'critical']
</script>

<template>
  <div class="lp">
    <button
      v-for="l in levels"
      :key="l"
      type="button"
      class="lvl"
      :class="[l, { on: model === l }]"
      @click="model = l"
    >
      <span class="dot" />
      {{ t(`level.${l}`) }}
    </button>
  </div>
</template>

<style lang="scss" scoped>
.lp {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}
.lvl {
  --c: var(--ink-faint);
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-family: var(--font-mono);
  font-size: 0.64rem;
  letter-spacing: 0.06em;
  padding: 7px 12px;
  border-radius: 9px;
  border: 1px solid var(--line);
  background: var(--surface);
  color: var(--ink-dim);
  cursor: pointer;
  transition: 0.2s;

  .dot { width: 8px; height: 8px; border-radius: 50%; background: var(--c); box-shadow: 0 0 8px var(--c); }
  &.info { --c: var(--primary); }
  &.notice { --c: var(--secondary); }
  &.warning { --c: var(--gold); }
  &.critical { --c: var(--ember); }

  &:hover { border-color: var(--line-strong); color: var(--ink); }
  &.on {
    color: var(--c);
    border-color: color-mix(in srgb, var(--c) 55%, transparent);
    background: color-mix(in srgb, var(--c) 12%, transparent);
  }
}
</style>
