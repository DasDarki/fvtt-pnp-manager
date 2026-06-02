<script setup lang="ts">
withDefaults(
  defineProps<{ label: string; value: string; percent: number; variant?: 'hp' | 'mana' }>(),
  { variant: 'mana' },
)
</script>

<template>
  <div class="aw-meter">
    <div class="ml">
      <span>{{ label }}</span>
      <span>{{ value }}</span>
    </div>
    <div class="mt">
      <i :class="variant" :style="{ width: `${Math.max(0, Math.min(100, percent))}%` }" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.aw-meter {
  .ml {
    display: flex;
    justify-content: space-between;
    font-family: var(--font-mono);
    font-size: 0.66rem;
    color: var(--ink-dim);
    margin-bottom: 6px;
  }
  .mt {
    height: 8px;
    border-radius: 999px;
    background: var(--surface-solid);
    overflow: hidden;
    border: 1px solid var(--line);

    i {
      display: block;
      height: 100%;
      border-radius: 999px;
      transition: width 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
    }
    i.hp {
      background: var(--grad-ember);
      box-shadow: 0 0 14px var(--ember);
    }
    i.mana {
      background: linear-gradient(90deg, var(--primary), var(--secondary));
      box-shadow: var(--glow-primary);
    }
  }
}
</style>
