<script setup lang="ts">
defineProps<{
  title: string
  icon: string
  count: number
  countLabel: string
  searchPlaceholder: string
}>()

const search = defineModel<string>('search', { default: '' })
</script>

<template>
  <header class="page-head">
    <div class="lead">
      <span class="ic"><Icon :name="icon" /></span>
      <div>
        <h1>{{ title }}</h1>
        <p>{{ count }} {{ countLabel }}</p>
      </div>
    </div>

    <div class="tools">
      <div class="search">
        <Icon name="lucide:search" />
        <input v-model="search" :placeholder="searchPlaceholder" />
      </div>
      <slot name="actions" />
    </div>
  </header>
</template>

<style lang="scss" scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin: 14px 0 26px;
  flex-wrap: wrap;
}

.lead {
  display: flex;
  align-items: center;
  gap: 16px;

  .ic {
    width: 50px;
    height: 50px;
    border-radius: 15px;
    display: grid;
    place-items: center;
    flex: none;
    background: color-mix(in srgb, var(--primary) 14%, transparent);
    color: var(--primary);
    border: 1px solid var(--line);
    box-shadow: var(--glow-primary);

    :deep(svg) { width: 24px; height: 24px; }
  }
  h1 {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: clamp(1.6rem, 3vw, 2.2rem);
    letter-spacing: 0.02em;
    line-height: 1.05;
    text-shadow: var(--glow-text);
  }
  p { color: var(--ink-faint); font-size: 0.82rem; font-family: var(--font-mono); margin-top: 4px; }
}

.tools {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search {
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 12px;
  padding: 10px 14px;
  transition: 0.25s;
  min-width: 220px;

  &:focus-within { border-color: var(--primary); box-shadow: var(--glow-primary); }
  :deep(svg) { width: 17px; height: 17px; color: var(--ink-faint); flex: none; }
  input {
    flex: 1;
    background: transparent;
    border: 0;
    outline: 0;
    color: var(--ink);
    font-family: var(--font-body);
    font-size: 0.9rem;
    min-width: 0;
  }
}

@media (max-width: 560px) {
  .tools { width: 100%; }
  .search { flex: 1; min-width: 0; }
}
</style>
