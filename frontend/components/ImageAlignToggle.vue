<script setup lang="ts">
const align = defineModel<string>({ default: 'center' })
const { t } = useI18n()
</script>

<template>
  <div class="ial">
    <span class="rl">{{ t('editor.align') }}</span>
    <div class="seg">
      <button
        v-for="a in IMAGE_ALIGNS"
        :key="a"
        type="button"
        :class="{ on: align === a }"
        :title="t(`align.${a}`)"
        @click="align = a"
      >
        <Icon :name="alignIcon(a)" />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.ial {
  display: flex;
  align-items: center;
  gap: 10px;
  .rl {
    font-family: var(--font-mono);
    font-size: 0.56rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: var(--ink-faint);
  }
  .seg {
    display: flex;
    margin-left: auto;
    border: 1px solid var(--line-strong);
    border-radius: 10px;
    overflow: hidden;
    background: var(--surface-2);
  }
  button {
    display: grid;
    place-items: center;
    width: 34px;
    height: 30px;
    border: 0;
    background: transparent;
    color: var(--ink-faint);
    cursor: pointer;
    transition: 0.18s;
    border-right: 1px solid var(--line);
    &:last-child { border-right: 0; }
    :deep(svg) { width: 15px; height: 15px; }
    &:hover { color: var(--ink); }
    &.on { color: var(--primary); background: rgba(70, 232, 255, 0.12); }
  }
}
</style>
