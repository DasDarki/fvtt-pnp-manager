<script setup lang="ts">
import type { Field, SystemData } from '~/types/ruleset'

const props = defineProps<{ field: Field; data: SystemData }>()

function abilityMod(v: any): string {
  const m = Math.floor((Number(v) - 10) / 2)
  return (m >= 0 ? '+' : '') + m
}

function setDot(trait: string, n: number) {
  const current = props.data[props.field.key]?.[trait]
  props.data[props.field.key][trait] = current === n ? n - 1 : n
}

const abShort: Record<string, string> = { str: 'STR', dex: 'GES', con: 'KON', int: 'INT', wis: 'WIS', cha: 'CHA' }

function profBonusVal(): number {
  const v = Number(props.data?.proficiency)
  return Number.isFinite(v) && v > 0 ? v : 2
}
function abilityScoreMod(ab?: string): number {
  if (!ab) return 0
  const score = Number(props.data?.abilities?.[ab] ?? 10)
  return Math.floor((score - 10) / 2)
}
function profTotal(trait: { key: string; ability?: string }): string {
  const lvl = props.data[props.field.key]?.[trait.key] ?? 0
  const add = lvl === 2 ? profBonusVal() * 2 : lvl === 1 ? profBonusVal() : 0
  const tot = abilityScoreMod(trait.ability) + add
  return (tot >= 0 ? '+' : '') + tot
}
function cycleProf(traitKey: string) {
  const max = props.field.proficiencyMax ?? 2
  const cur = props.data[props.field.key][traitKey] ?? 0
  props.data[props.field.key][traitKey] = cur >= max ? 0 : cur + 1
}
</script>

<template>
  <div class="sf" :class="`t-${field.type}`">
    <label v-if="field.type !== 'toggle'" class="lbl">{{ field.label }}</label>

    <input
      v-if="field.type === 'text'"
      v-model="data[field.key]"
      class="inp"
      type="text"
      :placeholder="field.placeholder"
    />

    <input
      v-else-if="field.type === 'number'"
      v-model.number="data[field.key]"
      class="inp"
      type="number"
      :min="field.min"
      :max="field.max"
    />

    <textarea
      v-else-if="field.type === 'textarea'"
      v-model="data[field.key]"
      class="inp"
      rows="4"
      :placeholder="field.placeholder"
    />

    <select v-else-if="field.type === 'select'" v-model="data[field.key]" class="inp">
      <option v-for="o in field.options" :key="o.value" :value="o.value">{{ o.label }}</option>
    </select>

    <label v-else-if="field.type === 'toggle'" class="toggle">
      <input v-model="data[field.key]" type="checkbox" />
      <span class="track" />
      <span class="thumb" />
      <span class="tl">{{ field.label }}</span>
    </label>

    <div v-else-if="field.type === 'abilities'" class="abilities">
      <div v-for="t in field.traits" :key="t.key" class="ab">
        <small>{{ t.label }}</small>
        <input v-model.number="data[field.key][t.key]" type="number" min="1" max="30" />
        <em>{{ abilityMod(data[field.key][t.key]) }}</em>
      </div>
    </div>

    <div v-else-if="field.type === 'dots'" class="dots-list">
      <div v-for="t in field.traits" :key="t.key" class="dot-row">
        <span class="dn">{{ t.label }}</span>
        <div class="dots">
          <button
            v-for="i in field.max || 5"
            :key="i"
            type="button"
            class="dot"
            :class="{ on: i <= data[field.key][t.key] }"
            :aria-label="`${t.label} ${i}`"
            @click="setDot(t.key, i)"
          />
        </div>
      </div>
    </div>

    <div v-else-if="field.type === 'proficiencies'" class="profs">
      <button
        v-for="t in field.traits"
        :key="t.key"
        type="button"
        class="prow"
        :class="`l${data[field.key][t.key]}`"
        @click="cycleProf(t.key)"
      >
        <span class="pin" />
        <span class="pn">{{ t.label }}<em v-if="t.ability">{{ abShort[t.ability] || t.ability.toUpperCase() }}</em></span>
        <span v-if="t.ability" class="pt">{{ profTotal(t) }}</span>
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.sf {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.t-abilities,
.t-dots {
  grid-column: 1 / -1;
}

.lbl {
  font-family: var(--font-mono);
  font-size: 0.62rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--ink-faint);
}

.inp {
  font-family: var(--font-body);
  font-size: 0.92rem;
  color: var(--ink);
  background: var(--surface-2);
  border: 1px solid var(--line-strong);
  border-radius: 11px;
  padding: 11px 13px;
  width: 100%;
  transition: 0.25s;

  &:focus {
    outline: 0;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(70, 232, 255, 0.14), var(--glow-primary);
  }
}
textarea.inp { resize: vertical; min-height: 90px; }
select.inp { cursor: pointer; appearance: none; }

.toggle {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding-left: 52px;
  min-height: 28px;

  input { display: none; }
  .track {
    position: absolute;
    left: 0;
    top: 0;
    width: 46px;
    height: 26px;
    border-radius: 999px;
    background: var(--surface-solid);
    border: 1px solid var(--line-strong);
    transition: 0.3s;
  }
  .thumb {
    position: absolute;
    left: 3px;
    top: 3px;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: var(--ink-dim);
    transition: 0.3s;
  }
  input:checked ~ .track { background: var(--grad-arcane); border-color: transparent; box-shadow: var(--glow-primary); }
  input:checked ~ .thumb { transform: translateX(20px); background: #fff; }
  .tl { font-size: 0.9rem; color: var(--ink-dim); }
}

.abilities {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 10px;

  .ab {
    text-align: center;
    background: var(--surface-2);
    border: 1px solid var(--line);
    border-radius: 12px;
    padding: 10px 4px;
    transition: 0.25s;

    &:focus-within { border-color: var(--primary); box-shadow: var(--glow-primary); }
    small { font-family: var(--font-mono); font-size: 0.58rem; letter-spacing: 0.1em; color: var(--ink-faint); display: block; }
    input {
      width: 100%;
      text-align: center;
      background: transparent;
      border: 0;
      outline: 0;
      color: var(--ink);
      font-family: var(--font-display);
      font-size: 1.5rem;
      font-weight: 600;
      -moz-appearance: textfield;
      &::-webkit-outer-spin-button,
      &::-webkit-inner-spin-button { -webkit-appearance: none; margin: 0; }
    }
    em { font-family: var(--font-mono); font-size: 0.74rem; font-style: normal; color: var(--primary); }
  }
}

.dots-list {
  display: flex;
  flex-direction: column;
  gap: 9px;

  .dot-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    background: var(--surface-2);
    border: 1px solid var(--line);
    border-radius: 11px;
    padding: 9px 14px;

    .dn { font-size: 0.88rem; color: var(--ink-dim); }
    .dots { display: flex; gap: 6px; }
    .dot {
      width: 16px;
      height: 16px;
      border-radius: 50%;
      border: 1px solid var(--line-strong);
      background: var(--surface-solid);
      cursor: pointer;
      transition: 0.18s;
      padding: 0;

      &:hover { border-color: var(--primary); }
      &.on { background: var(--grad-arcane); border-color: transparent; box-shadow: var(--glow-primary); }
    }
  }
}

.profs {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;

  .prow {
    display: flex;
    align-items: center;
    gap: 11px;
    padding: 9px 13px;
    border-radius: 11px;
    border: 1px solid var(--line);
    background: var(--surface-2);
    cursor: pointer;
    transition: 0.18s;
    text-align: left;
    color: var(--ink-dim);

    &:hover { border-color: var(--line-strong); }
    .pin {
      width: 14px;
      height: 14px;
      border-radius: 50%;
      flex: none;
      border: 1px solid var(--line-strong);
      background: var(--surface-solid);
      transition: 0.18s;
    }
    .pn {
      flex: 1;
      min-width: 0;
      font-size: 0.82rem;
      display: flex;
      align-items: baseline;
      gap: 6px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;

      em { font-family: var(--font-mono); font-size: 0.54rem; font-style: normal; letter-spacing: 0.08em; color: var(--ink-faint); }
    }
    .pt { font-family: var(--font-mono); font-size: 0.82rem; color: var(--ink-faint); flex: none; }

    &.l1 {
      color: var(--ink);
      border-color: color-mix(in srgb, var(--primary) 45%, transparent);
      .pin { background: var(--grad-arcane); border-color: transparent; box-shadow: var(--glow-primary); }
      .pt { color: var(--primary); }
    }
    &.l2 {
      color: var(--ink);
      border-color: color-mix(in srgb, var(--gold) 50%, transparent);
      .pin { background: var(--gold); border-color: transparent; box-shadow: 0 0 10px var(--gold); }
      .pt { color: var(--gold); }
    }
  }
}

@media (max-width: 560px) {
  .abilities { grid-template-columns: repeat(3, 1fr); }
  .profs { grid-template-columns: 1fr; }
}
</style>
