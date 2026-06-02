<script setup lang="ts">
import type { StatTile } from '~/types/entities'

const props = defineProps<{ stat: StatTile }>()
const { t } = useI18n()

const accentVar = computed(() => `var(--${props.stat.accent})`)
</script>

<template>
  <button class="stat" :style="{ '--c': accentVar }">
    <span v-if="stat.delta" class="delta" :class="stat.deltaTone">
      <Icon v-if="stat.deltaTone === 'up'" name="lucide:trending-up" />
      {{ stat.delta }}<template v-if="stat.deltaTone === 'flat'"> {{ t('stats.new') }}</template>
    </span>
    <span class="si"><Icon :name="stat.icon" /></span>
    <span class="big">{{ stat.value }}</span>
    <span class="lbl">{{ t(`stats.${stat.key}`) }}</span>
  </button>
</template>

<style lang="scss" scoped>
.stat {
  position: relative;
  padding: 20px;
  border-radius: 22px;
  background: var(--surface);
  border: 1px solid var(--line);
  box-shadow: var(--shadow-panel);
  overflow: hidden;
  transition: transform 0.3s, border-color 0.3s, box-shadow 0.3s;
  cursor: pointer;
  text-align: left;
  display: flex;
  flex-direction: column;
  align-items: flex-start;

  &:hover {
    transform: translateY(-4px);
    border-color: var(--line-strong);
    box-shadow: 0 28px 56px -28px #000, var(--glow-primary);
  }
  &::after {
    content: '';
    position: absolute;
    right: -30px;
    top: -30px;
    width: 110px;
    height: 110px;
    border-radius: 50%;
    background: var(--c);
    opacity: 0.14;
    filter: blur(8px);
  }

  .si {
    width: 38px;
    height: 38px;
    border-radius: 11px;
    display: grid;
    place-items: center;
    margin-bottom: 16px;
    background: color-mix(in srgb, var(--c) 15%, transparent);
    color: var(--c);

    :deep(svg) { width: 19px; height: 19px; }
  }
  .big {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: 2.1rem;
    line-height: 1;
    letter-spacing: 0.02em;
  }
  .lbl { font-size: 0.82rem; color: var(--ink-dim); margin-top: 6px; }
  .delta {
    position: absolute;
    top: 20px;
    right: 20px;
    font-family: var(--font-mono);
    font-size: 0.66rem;
    color: var(--emerald);
    display: flex;
    align-items: center;
    gap: 4px;

    &.flat { color: var(--gold); }
    :deep(svg) { width: 12px; height: 12px; }
  }
}
</style>
