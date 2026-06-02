<script setup lang="ts">
import type { Memory } from '~/types/entities'

const props = defineProps<{ memory: Memory; last?: boolean }>()

const levelColor: Record<string, string> = {
  info: 'var(--primary)',
  notice: 'var(--secondary)',
  warning: 'var(--gold)',
  critical: 'var(--ember)',
}
const subjectIcon: Record<string, string> = {
  character: 'lucide:user',
  item: 'lucide:gem',
  scene: 'lucide:castle',
  image: 'lucide:image',
  campaign: 'lucide:scroll-text',
}

const color = computed(() => levelColor[props.memory.level])
const icon = computed(() =>
  props.memory.level === 'critical' ? 'lucide:circle-alert' : subjectIcon[props.memory.subjectType],
)
</script>

<template>
  <div class="fitem" :class="{ line: !last }" :style="{ '--c': color }">
    <div class="stem">
      <span class="node"><Icon :name="icon" /></span>
    </div>
    <div class="ft">
      <b>{{ memory.title }}</b>
      <p>{{ memory.body }}</p>
      <span class="lk">↳ {{ memory.subjectLabel }}</span>
    </div>
    <span class="time">{{ memory.time }}</span>
  </div>
</template>

<style lang="scss" scoped>
.fitem {
  display: flex;
  gap: 15px;
  padding: 16px 0;
  border-bottom: 1px solid var(--line);
  position: relative;

  &:last-child { border-bottom: 0; }

  .stem {
    display: flex;
    flex-direction: column;
    align-items: center;
    flex: none;

    .node {
      width: 30px;
      height: 30px;
      border-radius: 9px;
      display: grid;
      place-items: center;
      background: color-mix(in srgb, var(--c) 16%, transparent);
      color: var(--c);
      border: 1px solid color-mix(in srgb, var(--c) 40%, transparent);

      :deep(svg) { width: 14px; height: 14px; }
    }
  }
  &.line .stem::after {
    content: '';
    flex: 1;
    width: 1px;
    background: var(--line);
    margin-top: 8px;
  }

  .ft {
    min-width: 0;
    b { font-size: 0.86rem; font-weight: 600; }
    p { font-size: 0.8rem; color: var(--ink-dim); margin-top: 3px; }
    .lk { font-family: var(--font-mono); font-size: 0.64rem; color: var(--c); margin-top: 7px; display: inline-block; }
  }
  .time {
    margin-left: auto;
    font-family: var(--font-mono);
    font-size: 0.6rem;
    color: var(--ink-faint);
    white-space: nowrap;
  }
}
</style>
