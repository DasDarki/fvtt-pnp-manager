import type { FieldGroup, RulesetAdapter, SystemData } from '~/types/ruleset'
import { dnd5e2024 } from '~/rulesets/dnd5e2024'
import { vampireV5 } from '~/rulesets/vampireV5'

const ADAPTERS: Record<string, RulesetAdapter> = {
  dnd5e_2024: dnd5e2024,
  vampire_v5: vampireV5,
}

export function getAdapter(id?: string | null): RulesetAdapter {
  return (id && ADAPTERS[id]) || dnd5e2024
}

export function buildDefaults(schema: FieldGroup[]): SystemData {
  const out: SystemData = {}
  for (const group of schema) {
    for (const f of group.fields) {
      if (f.type === 'number') out[f.key] = 0
      else if (f.type === 'toggle') out[f.key] = false
      else if (f.type === 'select') out[f.key] = f.options?.[0]?.value ?? ''
      else if (f.type === 'abilities') {
        out[f.key] = {}
        for (const t of f.traits ?? []) out[f.key][t.key] = 10
      } else if (f.type === 'dots') {
        out[f.key] = {}
        for (const t of f.traits ?? []) out[f.key][t.key] = 1
      } else if (f.type === 'proficiencies') {
        out[f.key] = {}
        for (const t of f.traits ?? []) out[f.key][t.key] = 0
      } else out[f.key] = ''
    }
  }
  return out
}

export function mergeSystemData(schema: FieldGroup[], existing?: SystemData | null): SystemData {
  const base = buildDefaults(schema)
  const sd = existing ?? {}
  for (const k of Object.keys(sd)) {
    const bv = base[k]
    const sv = sd[k]
    if (bv && typeof bv === 'object' && !Array.isArray(bv) && sv && typeof sv === 'object' && !Array.isArray(sv)) {
      base[k] = { ...bv, ...sv }
    } else {
      base[k] = sv
    }
  }
  return base
}
