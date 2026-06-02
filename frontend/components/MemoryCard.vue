<script setup lang="ts">
import type { Memory } from '~/types/entities'

const props = defineProps<{
  memory: Memory
  onAck?: () => void | Promise<void>
  onPin?: () => void | Promise<void>
  onDelete?: () => void | Promise<void>
}>()

const { t } = useI18n()

const levelColor: Record<string, string> = {
  info: 'var(--primary)',
  notice: 'var(--secondary)',
  warning: 'var(--gold)',
  critical: 'var(--ember)',
}
const color = computed(() => levelColor[props.memory.level])
const confirming = ref(false)

function del() {
  if (!confirming.value) {
    confirming.value = true
    setTimeout(() => (confirming.value = false), 3000)
    return
  }
  props.onDelete?.()
}
</script>

<template>
  <article class="mem" :class="{ ack: memory.acknowledged }" :style="{ '--c': color }">
    <span class="bar" />
    <div class="body">
      <div class="top">
        <span class="badge">{{ t(`level.${memory.level}`) }}</span>
        <b>{{ memory.title }}</b>
        <Icon v-if="memory.pinned" name="lucide:pin" class="pinned" />
        <span class="time">{{ memory.time }}</span>
      </div>
      <p v-if="memory.body">{{ memory.body }}</p>
      <span v-if="memory.subjectLabel" class="lk">↳ {{ memory.subjectLabel }}</span>
    </div>

    <div class="acts">
      <button
        v-if="memory.level === 'critical' && !memory.acknowledged"
        class="act ack-btn"
        :title="t('memory.ack')"
        @click="onAck?.()"
      >
        <Icon name="lucide:check" />
      </button>
      <span v-else-if="memory.acknowledged" class="done"><Icon name="lucide:check-check" /></span>
      <button class="act" :class="{ on: memory.pinned }" :title="t('memory.pin')" @click="onPin?.()">
        <Icon name="lucide:pin" />
      </button>
      <button class="act danger" :title="t('memory.delete')" @click="del">
        <Icon :name="confirming ? 'lucide:trash-2' : 'lucide:trash'" />
      </button>
    </div>
  </article>
</template>

<style lang="scss" scoped>
.mem {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 16px 18px 16px 20px;
  border-radius: 16px;
  border: 1px solid var(--line);
  background: var(--surface);
  box-shadow: var(--shadow-panel);
  overflow: hidden;
  transition: 0.25s;

  &:hover { border-color: var(--line-strong); }
  &.ack { opacity: 0.62; }

  .bar { position: absolute; left: 0; top: 0; bottom: 0; width: 3px; background: var(--c); box-shadow: 0 0 12px var(--c); }

  .body { flex: 1; min-width: 0; }
  .top {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;

    .badge {
      font-family: var(--font-mono);
      font-size: 0.54rem;
      letter-spacing: 0.12em;
      text-transform: uppercase;
      color: var(--c);
      border: 1px solid color-mix(in srgb, var(--c) 50%, transparent);
      border-radius: 999px;
      padding: 2px 9px;
    }
    b { font-family: var(--font-display); font-weight: 600; font-size: 0.98rem; }
    .pinned { width: 13px; height: 13px; color: var(--gold); }
    .time { margin-left: auto; font-family: var(--font-mono); font-size: 0.62rem; color: var(--ink-faint); white-space: nowrap; }
  }
  p { font-size: 0.84rem; color: var(--ink-dim); margin-top: 6px; }
  .lk { display: inline-block; margin-top: 8px; font-family: var(--font-mono); font-size: 0.64rem; color: var(--c); }

  .acts { display: flex; gap: 6px; flex: none; }
  .act {
    width: 30px;
    height: 30px;
    display: grid;
    place-items: center;
    border-radius: 8px;
    border: 1px solid var(--line);
    background: var(--surface-2);
    color: var(--ink-faint);
    cursor: pointer;
    transition: 0.2s;
    :deep(svg) { width: 15px; height: 15px; }
    &:hover { color: var(--ink); border-color: var(--line-strong); }
    &.on { color: var(--gold); border-color: var(--gold); }
    &.ack-btn:hover { color: var(--emerald); border-color: var(--emerald); }
    &.danger:hover { color: var(--ember); border-color: var(--ember); }
  }
  .done { width: 30px; height: 30px; display: grid; place-items: center; color: var(--emerald); :deep(svg) { width: 16px; height: 16px; } }
}
</style>
