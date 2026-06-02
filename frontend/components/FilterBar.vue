<script setup lang="ts">
import type { ApiTag } from '~/types/api'

withDefaults(
  defineProps<{
    primaryOptions?: { value: string; label: string }[]
    sortOptions: { value: string; label: string }[]
    campaignTags?: ApiTag[]
    showImageFilter?: boolean
  }>(),
  { showImageFilter: false },
)

const { t } = useI18n()
const primary = defineModel<string>('primary')
const tags = defineModel<string[]>('tags', { default: () => [] })
const sort = defineModel<string>('sort')
const hasImage = defineModel<boolean>('hasImage')

function toggleTag(id: string) {
  const a = [...(tags.value || [])]
  const i = a.indexOf(id)
  if (i >= 0) a.splice(i, 1)
  else a.push(id)
  tags.value = a
}
</script>

<template>
  <div class="fb">
    <div v-if="primaryOptions?.length" class="chips">
      <button
        v-for="o in primaryOptions"
        :key="o.value"
        class="chip"
        :class="{ on: primary === o.value }"
        @click="primary = o.value"
      >
        {{ o.label }}
      </button>
    </div>

    <div v-if="campaignTags?.length" class="chips tags">
      <button
        v-for="tg in campaignTags"
        :key="tg.id"
        class="chip tag"
        :class="{ on: tags?.includes(tg.id) }"
        :style="{ '--c': `var(--${tg.color})` }"
        @click="toggleTag(tg.id)"
      >
        <span class="d" /> {{ tg.name }}
      </button>
    </div>

    <div class="right">
      <label v-if="showImageFilter" class="imgf" :class="{ on: hasImage }">
        <input v-model="hasImage" type="checkbox" />
        <Icon name="lucide:image" /> {{ t('filterbar.hasImage') }}
      </label>
      <div class="sortw">
        <Icon name="lucide:arrow-down-up" />
        <select v-model="sort" class="sort">
          <option v-for="s in sortOptions" :key="s.value" :value="s.value">{{ s.label }}</option>
        </select>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.fb {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
  margin-bottom: 22px;
}

.chips { display: flex; gap: 8px; flex-wrap: wrap; }
.chips.tags { padding-left: 14px; border-left: 1px solid var(--line); }

.chip {
  font-family: var(--font-mono);
  font-size: 0.66rem;
  letter-spacing: 0.06em;
  color: var(--ink-dim);
  padding: 7px 13px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: var(--surface);
  cursor: pointer;
  transition: 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 6px;

  &:hover { color: var(--ink); border-color: var(--line-strong); }
  &.on { color: #06040c; background: var(--grad-arcane); border-color: transparent; box-shadow: var(--glow-primary); }

  &.tag {
    .d { width: 7px; height: 7px; border-radius: 50%; background: var(--c, var(--secondary)); box-shadow: 0 0 6px var(--c, var(--secondary)); }
    &.on {
      color: var(--c, var(--secondary));
      background: color-mix(in srgb, var(--c, var(--secondary)) 14%, transparent);
      border-color: color-mix(in srgb, var(--c, var(--secondary)) 55%, transparent);
      box-shadow: 0 0 14px -6px var(--c, var(--secondary));
    }
  }
}
:global(html[data-theme='light']) .chip.on:not(.tag) { color: #fff; }

.right { margin-left: auto; display: flex; align-items: center; gap: 10px; }

.imgf {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-family: var(--font-mono);
  font-size: 0.66rem;
  color: var(--ink-dim);
  border: 1px solid var(--line);
  background: var(--surface);
  border-radius: 999px;
  padding: 7px 13px;
  cursor: pointer;
  transition: 0.2s;
  input { display: none; }
  :deep(svg) { width: 14px; height: 14px; }
  &:hover { color: var(--ink); border-color: var(--line-strong); }
  &.on { color: var(--primary); border-color: var(--primary); box-shadow: var(--glow-primary); }
}

.sortw {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  border: 1px solid var(--line);
  background: var(--surface);
  border-radius: 999px;
  padding: 4px 12px 4px 13px;
  :deep(svg) { width: 14px; height: 14px; color: var(--ink-faint); }
  .sort {
    appearance: none;
    border: 0;
    background: transparent;
    color: var(--ink-dim);
    font-family: var(--font-mono);
    font-size: 0.66rem;
    cursor: pointer;
    padding: 5px 2px;
    outline: 0;
  }
}
</style>
