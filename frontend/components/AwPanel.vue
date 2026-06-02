<script setup lang="ts">
withDefaults(defineProps<{ pad?: boolean }>(), { pad: false })
</script>

<template>
  <section class="aw-panel" :class="{ pad }">
    <slot />
  </section>
</template>

<style lang="scss" scoped>
.aw-panel {
  position: relative;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 22px;
  backdrop-filter: blur(16px);
  box-shadow: var(--shadow-panel);
  overflow: hidden;

  &.pad {
    padding: 26px;
  }
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--line-strong), transparent);
  }
}
</style>
