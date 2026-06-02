<script setup lang="ts">
import type { CharacterStatus } from '~/types/entities'

defineProps<{ status: CharacterStatus | 'prepared' }>()
</script>

<template>
  <span class="aw-pill" :class="status">
    <span class="led" />
    <slot />
  </span>
</template>

<style lang="scss" scoped>
.aw-pill {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-family: var(--font-mono);
  font-size: 0.6rem;
  letter-spacing: 0.1em;
  font-weight: 500;
  padding: 5px 11px;
  border-radius: 999px;
  border: 1px solid currentColor;
  text-transform: uppercase;

  .led {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: currentColor;
    box-shadow: 0 0 8px currentColor;
  }

  &.alive { color: var(--emerald); }
  &.dead { color: var(--ember); }
  &.unknown { color: var(--ink-faint); }
  &.hunted,
  &.prepared { color: var(--gold); }
}
</style>
