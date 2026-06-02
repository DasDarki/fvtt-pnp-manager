<script setup lang="ts">
const { t, locale, setLocale } = useI18n()
const { theme, toggle } = useTheme()
const ui = useUiStore()
</script>

<template>
  <header class="tbar">
    <button class="cmd" @click="ui.openPalette()">
      <Icon name="lucide:search" />
      <span>{{ t('search.placeholder') }}</span>
      <span class="kbd">Ctrl K</span>
    </button>

    <div class="sp" />

    <div class="seg hide-m" role="group">
      <button :class="{ on: locale === 'de' }" @click="setLocale('de')">DE</button>
      <button :class="{ on: locale === 'en' }" @click="setLocale('en')">EN</button>
    </div>

    <button class="iconbtn" :aria-label="t('a11y.toggleTheme')" @click="toggle">
      <Icon :name="theme === 'dark' ? 'lucide:moon-star' : 'lucide:sun'" />
    </button>

    <AwButton icon="lucide:plus" @click="ui.openPalette()">
      <span class="hide-m">{{ t('actions.create') }}</span>
    </AwButton>
  </header>
</template>

<style lang="scss" scoped>
.tbar {
  position: sticky;
  top: 0;
  z-index: 30;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 30px;
  background: linear-gradient(180deg, var(--void) 40%, transparent);
  backdrop-filter: blur(12px);
}

.cmd {
  flex: 1;
  max-width: 480px;
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 13px;
  padding: 11px 15px;
  cursor: text;
  transition: 0.25s;
  color: var(--ink-faint);
  font-family: var(--font-body);

  :deep(svg) { width: 17px; height: 17px; }
  span { flex: 1; text-align: left; font-size: 0.9rem; }
  .kbd {
    flex: none;
    font-family: var(--font-mono);
    font-size: 0.64rem;
    color: var(--ink-dim);
    border: 1px solid var(--line-strong);
    border-radius: 6px;
    padding: 3px 7px;
    background: var(--surface);
  }
  &:hover { border-color: var(--line-strong); box-shadow: var(--glow-primary); }
}

.sp { flex: 1; }

.seg {
  display: inline-flex;
  padding: 3px;
  gap: 2px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: var(--surface);

  button {
    font-family: var(--font-mono);
    font-size: 0.64rem;
    letter-spacing: 0.1em;
    font-weight: 500;
    border: 0;
    background: transparent;
    color: var(--ink-faint);
    padding: 6px 11px;
    border-radius: 999px;
    cursor: pointer;
    transition: 0.25s;

    &.on {
      background: var(--grad-arcane);
      color: #06040c;
      box-shadow: var(--glow-primary);
    }
  }
}
:global(html[data-theme='light']) .seg button.on { color: #fff; }

.iconbtn {
  width: 40px;
  height: 40px;
  display: grid;
  place-items: center;
  border-radius: 12px;
  cursor: pointer;
  border: 1px solid var(--line-strong);
  background: var(--surface);
  color: var(--ink);
  transition: 0.25s;

  :deep(svg) { width: 18px; height: 18px; }
  &:hover { color: var(--primary); border-color: var(--primary); box-shadow: var(--glow-primary); transform: translateY(-1px); }
}

@media (max-width: 820px) {
  .tbar { padding-left: 70px; }
  .hide-m { display: none; }
}
</style>
