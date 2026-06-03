<script setup lang="ts">
import type { SceneSummary } from '~/types/entities'

defineProps<{ scene: SceneSummary }>()
</script>

<template>
  <article class="scene">
    <div class="bg" :class="`tone-${scene.tone || 'arcane'}`" aria-hidden="true" />
    <img v-if="scene.image" class="map" :src="scene.image" :alt="scene.title" :style="{ objectPosition: objPos(scene.imageAlign) }" />
    <div class="veil" aria-hidden="true" />
    <div class="body">
      <div class="tagrow">
        <AwRarity tier="act">{{ scene.act }}</AwRarity>
        <AwPill status="prepared">{{ scene.status }}</AwPill>
      </div>
      <h3>{{ scene.title }}</h3>
      <p>{{ scene.description }}</p>
      <TagChips v-if="scene.tags?.length" :tags="scene.tags" class="sc-tags" />
    </div>
    <div class="actors" aria-hidden="true">
      <span v-for="(a, i) in scene.actors" :key="i" class="a">{{ a }}</span>
      <span v-if="scene.extra" class="a more">+{{ scene.extra }}</span>
    </div>
  </article>
</template>

<style lang="scss" scoped>
.scene {
  position: relative;
  border-radius: 22px;
  overflow: hidden;
  min-height: 230px;
  display: flex;
  align-items: flex-end;
  padding: 24px;
  box-shadow: var(--shadow-panel);
  border: 1px solid var(--line);

  .bg {
    position: absolute;
    inset: 0;
    z-index: 0;
  }
  .bg.tone-arcane {
    background:
      radial-gradient(120% 90% at 72% 12%, rgba(183, 104, 255, 0.55), transparent 55%),
      radial-gradient(90% 80% at 20% 90%, rgba(70, 232, 255, 0.4), transparent 60%),
      linear-gradient(160deg, #1a1030, #0a1620);
  }
  .bg.tone-ember {
    background:
      radial-gradient(120% 90% at 70% 14%, rgba(255, 194, 77, 0.45), transparent 55%),
      radial-gradient(90% 80% at 20% 95%, rgba(255, 106, 85, 0.4), transparent 60%),
      linear-gradient(160deg, #2a1410, #1a0a14);
  }
  .bg.tone-verdant {
    background:
      radial-gradient(120% 90% at 72% 12%, rgba(55, 232, 164, 0.4), transparent 55%),
      radial-gradient(90% 80% at 22% 92%, rgba(70, 232, 255, 0.32), transparent 60%),
      linear-gradient(160deg, #0c2018, #0a1614);
  }
  .map {
    position: absolute;
    inset: 0;
    z-index: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .veil {
    position: absolute;
    inset: 0;
    z-index: 1;
    background: linear-gradient(to top, rgba(7, 6, 13, 0.94), transparent 65%);
  }
  .body { position: relative; z-index: 2; width: 100%; }
  .tagrow { display: flex; gap: 9px; margin-bottom: 12px; }
  h3 { font-family: var(--font-display); font-weight: 600; font-size: 1.5rem; letter-spacing: 0.02em; color: var(--ink); }
  p { font-size: 0.86rem; color: var(--ink-dim); margin-top: 7px; max-width: 62%; }
  .sc-tags { margin-top: 10px; }

  .actors {
    position: absolute;
    right: 24px;
    bottom: 24px;
    display: flex;
    z-index: 2;

    .a {
      width: 38px;
      height: 38px;
      border-radius: 11px;
      border: 2px solid var(--surface-solid);
      margin-left: -12px;
      display: grid;
      place-items: center;
      font-family: var(--font-display);
      font-size: 0.85rem;
      font-weight: 600;
      color: #06040c;
      background: var(--grad-arcane);

      &:nth-child(2) { background: linear-gradient(120deg, var(--gold), var(--ember)); }
      &:nth-child(3) { background: linear-gradient(120deg, var(--secondary), var(--magenta)); }
      &.more {
        background: var(--surface-solid);
        color: var(--ink);
        font-family: var(--font-mono);
        font-size: 0.68rem;
      }
    }
  }
}
</style>
