<script setup lang="ts">
import type { Rarity } from '~/types/entities'

defineProps<{ tier: Rarity | 'act' }>()
</script>

<template>
  <span class="aw-rarity" :class="tier">
    <slot />
  </span>
</template>

<style lang="scss" scoped>
.aw-rarity {
  --c: var(--ink-dim);
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-family: var(--font-mono);
  font-size: 0.62rem;
  letter-spacing: 0.06em;
  color: var(--c);
  padding: 5px 11px;
  border-radius: 8px;
  border: 1px solid color-mix(in srgb, var(--c) 45%, transparent);
  background: color-mix(in srgb, var(--c) 11%, transparent);

  &::before {
    content: '';
    width: 7px;
    height: 7px;
    transform: rotate(45deg);
    background: var(--c);
    box-shadow: 0 0 8px var(--c);
  }

  &.common { --c: var(--ink-faint); }
  &.uncommon { --c: var(--emerald); }
  &.rare { --c: var(--primary); }
  &.epic,
  &.act { --c: var(--secondary); }
  &.legendary { --c: var(--gold); }
  &.artifact { --c: var(--magenta); }
}
</style>
