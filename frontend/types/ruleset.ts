import type { MiniStat } from './entities'

export type FieldType =
  | 'text'
  | 'number'
  | 'textarea'
  | 'select'
  | 'toggle'
  | 'abilities'
  | 'dots'
  | 'proficiencies'

export interface FieldOption {
  value: string
  label: string
}

export interface Trait {
  key: string
  label: string
  ability?: string
}

export interface Field {
  key: string
  label: string
  type: FieldType
  options?: FieldOption[]
  traits?: Trait[]
  min?: number
  max?: number
  proficiencyMax?: number
  placeholder?: string
  hint?: string
}

export interface FieldGroup {
  title: string
  icon?: string
  fields: Field[]
}

export type SystemData = Record<string, any>

export interface RulesetAdapter {
  id: string
  label: string
  schema: FieldGroup[]
  subtitle(data: SystemData): string
  cardStats(data: SystemData): MiniStat[]
  hpPercent(data: SystemData): number
}
