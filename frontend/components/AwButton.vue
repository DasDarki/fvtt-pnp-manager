<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'primary' | 'ghost' | 'soft' | 'danger'
    icon?: string
    square?: boolean
    type?: 'button' | 'submit' | 'reset'
  }>(),
  { variant: 'primary', type: 'button' },
)
</script>

<template>
  <button :type="type" class="aw-btn" :class="[`v-${variant}`, { square }]">
    <Icon v-if="icon" :name="icon" class="ic" />
    <span v-if="$slots.default" class="lbl"><slot /></span>
  </button>
</template>

<style lang="scss" scoped>
.aw-btn {
  font-family: var(--font-body);
  font-weight: 600;
  font-size: 0.86rem;
  letter-spacing: 0.01em;
  border: 0;
  cursor: pointer;
  padding: 11px 20px;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  gap: 9px;
  transition: transform 0.2s, box-shadow 0.25s, background 0.25s, border-color 0.25s, color 0.25s;
  position: relative;
  white-space: nowrap;

  .ic {
    width: 16px;
    height: 16px;
  }
  &.square {
    padding: 11px;
  }
  &:disabled {
    opacity: 0.4;
    cursor: not-allowed;
    transform: none !important;
    box-shadow: none !important;
  }
}

.v-primary {
  background: var(--grad-arcane);
  color: #06040c;
  box-shadow: var(--glow-primary);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 0 36px rgba(70, 232, 255, 0.6), 0 10px 30px -8px rgba(183, 104, 255, 0.5);
  }
}
:global(html[data-theme='light']) .v-primary {
  color: #fff;
}

.v-ghost {
  background: var(--surface);
  color: var(--ink);
  border: 1px solid var(--line-strong);

  &:hover {
    transform: translateY(-2px);
    color: var(--primary);
    border-color: var(--primary);
    box-shadow: var(--glow-primary);
  }
}

.v-soft {
  background: rgba(183, 104, 255, 0.12);
  color: var(--secondary);
  border: 1px solid rgba(183, 104, 255, 0.32);

  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--glow-secondary);
  }
}

.v-danger {
  background: rgba(255, 106, 85, 0.12);
  color: var(--ember);
  border: 1px solid rgba(255, 106, 85, 0.34);

  &:hover {
    transform: translateY(-2px);
    background: var(--ember);
    color: #1a0a06;
    box-shadow: 0 0 26px rgba(255, 106, 85, 0.5);
  }
}
</style>
