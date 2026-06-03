<script setup lang="ts">
import type { ImageEntry } from '~/types/entities'

const props = defineProps<{ image: ImageEntry }>()
const { t } = useI18n()

const pushLabel = computed(() => (props.image.pushAs === 'journal' ? 'Journal' : 'Actor'))
</script>

<template>
  <article class="ic">
    <div class="thumb">
      <img v-if="image.image" :src="image.image" :alt="image.name" loading="lazy" :style="{ objectPosition: objPos(image.imageAlign) }" />
      <Icon v-else name="lucide:image" class="ph" />
      <span class="badge"><Icon :name="image.pushAs === 'journal' ? 'lucide:scroll-text' : 'lucide:user'" /> {{ pushLabel }}</span>
    </div>
    <div class="body">
      <h3>{{ image.name }}</h3>
      <p v-if="image.notes">{{ image.notes }}</p>
      <p v-else class="muted">{{ t('image.noNotes') }}</p>
      <TagChips v-if="image.tags?.length" :tags="image.tags" class="im-tags" />
    </div>
  </article>
</template>

<style lang="scss" scoped>
.ic {
  position: relative;
  border-radius: 20px;
  overflow: hidden;
  background: var(--surface);
  border: 1px solid var(--line);
  box-shadow: var(--shadow-panel);
  transition: transform 0.3s, box-shadow 0.3s, border-color 0.3s;

  &:hover {
    transform: translateY(-6px);
    border-color: var(--line-strong);
    box-shadow: 0 30px 60px -30px #000, var(--glow-secondary);
  }
}

.thumb {
  position: relative;
  aspect-ratio: 4 / 3;
  background: radial-gradient(120% 100% at 30% 20%, rgba(183, 104, 255, 0.18), transparent 60%), var(--void-2);
  display: grid;
  place-items: center;
  overflow: hidden;

  img { width: 100%; height: 100%; object-fit: cover; }
  .ph { width: 38px; height: 38px; color: var(--ink-faint); opacity: 0.5; }
  .badge {
    position: absolute;
    top: 10px;
    left: 10px;
    display: inline-flex;
    align-items: center;
    gap: 5px;
    font-family: var(--font-mono);
    font-size: 0.56rem;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--ink);
    background: rgba(7, 6, 13, 0.6);
    backdrop-filter: blur(6px);
    border: 1px solid var(--line-strong);
    padding: 4px 9px;
    border-radius: 999px;
    :deep(svg) { width: 11px; height: 11px; color: var(--secondary); }
  }
}

.body {
  padding: 14px 16px;
  h3 { font-family: var(--font-display); font-weight: 600; font-size: 1rem; line-height: 1.2; }
  p { font-size: 0.78rem; color: var(--ink-dim); margin-top: 5px; display: -webkit-box; -webkit-line-clamp: 2; line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
  .muted { color: var(--ink-faint); font-style: italic; }
  .im-tags { margin-top: 9px; }
}
</style>
