<script setup lang="ts">
import type { Item } from '~/types/entities'

const props = defineProps<{ item: Item }>()
const { t } = useI18n()

const rarityColor: Record<string, string> = {
  common: 'var(--ink-faint)',
  uncommon: 'var(--emerald)',
  rare: 'var(--primary)',
  epic: 'var(--secondary)',
  legendary: 'var(--gold)',
  artifact: 'var(--magenta)',
}
const accent = computed(() => rarityColor[props.item.rarity])
</script>

<template>
  <article class="item" :style="{ '--c': accent }">
    <div class="top">
      <span class="ibox">
        <img v-if="item.image" :src="item.image" :alt="item.name" />
        <Icon v-else :name="item.icon" />
      </span>
      <Icon v-if="item.attuned" name="lucide:link" class="attune" :title="t('item.attuned')" />
    </div>
    <div class="meta">
      <h3>{{ item.name }}</h3>
      <p>{{ item.type }}</p>
    </div>
    <div class="foot">
      <AwRarity :tier="item.rarity">{{ t(`rarity.${item.rarity}`) }}</AwRarity>
      <TagChips v-if="item.tags?.length" :tags="item.tags" :max="2" class="it-tags" />
    </div>
  </article>
</template>

<style lang="scss" scoped>
.item {
  position: relative;
  border-radius: 20px;
  padding: 18px;
  background: var(--surface);
  border: 1px solid var(--line);
  box-shadow: var(--shadow-panel);
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s, border-color 0.3s;
  display: flex;
  flex-direction: column;
  gap: 16px;

  &:hover {
    transform: translateY(-6px);
    border-color: color-mix(in srgb, var(--c) 55%, transparent);
    box-shadow: 0 30px 60px -30px #000, 0 0 26px -6px var(--c);
  }
  &::after {
    content: '';
    position: absolute;
    right: -26px;
    top: -26px;
    width: 96px;
    height: 96px;
    border-radius: 50%;
    background: var(--c);
    opacity: 0.12;
    filter: blur(6px);
  }

  .top {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
  }
  .ibox {
    width: 46px;
    height: 46px;
    border-radius: 13px;
    display: grid;
    place-items: center;
    background: color-mix(in srgb, var(--c) 16%, transparent);
    color: var(--c);
    border: 1px solid color-mix(in srgb, var(--c) 32%, transparent);
    box-shadow: 0 0 18px -6px var(--c);
    overflow: hidden;

    :deep(svg) { width: 22px; height: 22px; }
    img { width: 100%; height: 100%; object-fit: cover; border-radius: 12px; }
  }
  .attune { width: 15px; height: 15px; color: var(--c); opacity: 0.8; }

  .meta {
    h3 { font-family: var(--font-display); font-weight: 600; font-size: 1.02rem; line-height: 1.2; }
    p { font-size: 0.76rem; color: var(--ink-faint); margin-top: 4px; }
  }
  .foot { margin-top: auto; display: flex; flex-direction: column; gap: 9px; align-items: flex-start; }
}
</style>
