<script setup lang="ts">
import type { TagRef } from '~/types/entities'

withDefaults(defineProps<{ tags?: TagRef[]; max?: number }>(), { max: 4 })
</script>

<template>
  <div v-if="tags && tags.length" class="tagchips">
    <span v-for="tg in tags.slice(0, max)" :key="tg.id" class="tc" :style="{ '--c': `var(--${tg.color})` }">{{ tg.name }}</span>
    <span v-if="tags.length > max" class="more">+{{ tags.length - max }}</span>
  </div>
</template>

<style lang="scss" scoped>
.tagchips { display: flex; flex-wrap: wrap; gap: 5px; }
.tc {
  font-family: var(--font-mono);
  font-size: 0.54rem;
  letter-spacing: 0.04em;
  color: var(--c, var(--secondary));
  border: 1px solid color-mix(in srgb, var(--c, var(--secondary)) 40%, transparent);
  background: color-mix(in srgb, var(--c, var(--secondary)) 10%, transparent);
  padding: 2px 8px;
  border-radius: 999px;
}
.more { font-family: var(--font-mono); font-size: 0.54rem; color: var(--ink-faint); align-self: center; }
</style>
