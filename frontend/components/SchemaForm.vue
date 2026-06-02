<script setup lang="ts">
import type { FieldGroup, SystemData } from '~/types/ruleset'

defineProps<{ schema: FieldGroup[]; data: SystemData }>()
</script>

<template>
  <div class="schema-form">
    <AwPanel v-for="group in schema" :key="group.title" class="group">
      <div class="gh">
        <Icon v-if="group.icon" :name="group.icon" />
        <h3>{{ group.title }}</h3>
      </div>
      <div class="grid">
        <SchemaField v-for="f in group.fields" :key="f.key" :field="f" :data="data" />
      </div>
    </AwPanel>
  </div>
</template>

<style lang="scss" scoped>
.schema-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}
.group { padding: 22px; }

.gh {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 18px;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--line);

  :deep(svg) { width: 17px; height: 17px; color: var(--primary); }
  h3 { font-family: var(--font-display); font-weight: 600; font-size: 1.05rem; letter-spacing: 0.02em; }
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

@media (max-width: 620px) {
  .grid { grid-template-columns: 1fr; }
}
</style>
