<script setup lang="ts">
import type { Character } from '~/types/entities'

const props = defineProps<{ character: Character }>()
const { t } = useI18n()

const dotColor = computed(() => {
  const map: Record<string, string> = {
    alive: 'var(--emerald)',
    dead: 'var(--ember)',
    unknown: 'var(--ink-faint)',
    hunted: 'var(--gold)',
  }
  return map[props.character.status]
})
</script>

<template>
  <article class="card">
    <span v-if="character.critical" class="flag">
      <Icon name="lucide:triangle-alert" />
      {{ t('critical.badge') }}
    </span>

    <div class="cc-top">
      <div class="avatar" :style="{ background: character.ring }">
        <img v-if="character.image" :src="character.image" :alt="character.name" />
        <span v-else>{{ character.initial }}</span>
        <i class="sdot" :style="{ background: dotColor, boxShadow: `0 0 10px ${dotColor}` }" />
      </div>
      <div class="cc-id">
        <div class="cc-name">{{ character.name }}</div>
        <div class="cc-sub">{{ character.subtitle }}</div>
      </div>
    </div>

    <TagChips v-if="character.tags?.length" :tags="character.tags" class="cc-tags" />

    <div class="cc-stats">
      <div v-for="s in character.stats" :key="s.label" class="s">
        <b>{{ s.value }}</b>
        <small>{{ s.label }}</small>
      </div>
    </div>

    <div class="cc-bar"><i :style="{ width: `${character.hpPercent}%` }" /></div>
  </article>
</template>

<style lang="scss" scoped>
.card {
  position: relative;
  border-radius: 22px;
  overflow: hidden;
  background: var(--surface);
  border: 1px solid var(--line);
  box-shadow: var(--shadow-panel);
  transition: transform 0.3s, box-shadow 0.3s, border-color 0.3s;

  &:hover {
    transform: translateY(-6px);
    border-color: var(--line-strong);
    box-shadow: 0 30px 60px -30px #000, var(--glow-primary);
  }
  &::after {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    pointer-events: none;
    background: linear-gradient(160deg, rgba(255, 255, 255, 0.06), transparent 42%);
  }
}

.flag {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 3;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-family: var(--font-mono);
  font-size: 0.52rem;
  letter-spacing: 0.12em;
  font-weight: 600;
  color: #1a0a06;
  background: var(--grad-ember);
  padding: 4px 9px;
  border-radius: 999px;
  box-shadow: 0 0 16px -2px var(--ember);
  text-transform: uppercase;

  :deep(svg) { width: 11px; height: 11px; }
}

.cc-top {
  padding: 18px;
  display: flex;
  gap: 14px;
  align-items: center;
}

.avatar {
  width: 54px;
  height: 54px;
  border-radius: 14px;
  flex: none;
  position: relative;
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 600;
  font-size: 1.3rem;
  color: var(--ink);
  box-shadow: var(--glow-secondary);

  &::before {
    content: '';
    position: absolute;
    inset: 2px;
    border-radius: 12px;
    background: radial-gradient(circle at 30% 25%, var(--surface-2), var(--surface-solid));
  }
  span { position: relative; z-index: 1; }
  img {
    position: absolute;
    inset: 2px;
    border-radius: 12px;
    object-fit: cover;
    z-index: 1;
  }
  .sdot {
    position: absolute;
    bottom: -3px;
    right: -3px;
    width: 15px;
    height: 15px;
    border-radius: 50%;
    border: 3px solid var(--surface-solid);
    z-index: 2;
  }
}

.cc-id { min-width: 0; }
.cc-name { font-family: var(--font-display); font-weight: 600; font-size: 1rem; line-height: 1.2; }
.cc-sub { font-size: 0.74rem; color: var(--ink-faint); margin-top: 3px; }
.cc-tags { padding: 0 18px 13px; margin-top: -2px; }

.cc-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  border-top: 1px solid var(--line);

  .s {
    padding: 11px 6px;
    text-align: center;
    border-right: 1px solid var(--line);

    &:last-child { border-right: 0; }
    b {
      font-family: var(--font-mono);
      font-size: 1rem;
      color: var(--primary);
      display: block;
      text-shadow: var(--glow-text);
    }
    small {
      font-family: var(--font-mono);
      font-size: 0.54rem;
      letter-spacing: 0.12em;
      color: var(--ink-faint);
      text-transform: uppercase;
    }
  }
}

.cc-bar {
  height: 3px;
  background: var(--line);

  i {
    display: block;
    height: 100%;
    background: var(--grad-ember);
    box-shadow: 0 0 12px var(--ember);
    transition: width 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
  }
}
</style>
