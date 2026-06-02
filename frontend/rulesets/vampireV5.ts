import type { RulesetAdapter, SystemData } from '~/types/ruleset'

function n(data: SystemData, key: string): number {
  const v = Number(data?.[key])
  return Number.isFinite(v) ? v : 0
}

export const vampireV5: RulesetAdapter = {
  id: 'vampire_v5',
  label: 'Vampire V5',
  schema: [
    {
      title: 'Grundlagen',
      icon: 'lucide:scroll-text',
      fields: [
        {
          key: 'clan',
          label: 'Clan',
          type: 'select',
          options: [
            { value: 'brujah', label: 'Brujah' },
            { value: 'gangrel', label: 'Gangrel' },
            { value: 'malkavian', label: 'Malkavianer' },
            { value: 'nosferatu', label: 'Nosferatu' },
            { value: 'toreador', label: 'Toreador' },
            { value: 'tremere', label: 'Tremere' },
            { value: 'ventrue', label: 'Ventrue' },
            { value: 'lasombra', label: 'Lasombra' },
            { value: 'caitiff', label: 'Caitiff' },
          ],
        },
        { key: 'generation', label: 'Generation', type: 'number', min: 4, max: 16 },
        { key: 'predator', label: 'Raubtiertyp', type: 'text', placeholder: 'z.B. Alleycat' },
        { key: 'sire', label: 'Erzeuger', type: 'text' },
        { key: 'subtitle', label: 'Untertitel (Karte)', type: 'text', placeholder: 'Toreador · Gen 11' },
      ],
    },
    {
      title: 'Attribute',
      icon: 'lucide:dices',
      fields: [
        {
          key: 'attributes',
          label: 'Attribute (1–5)',
          type: 'dots',
          max: 5,
          traits: [
            { key: 'strength', label: 'Stärke' },
            { key: 'dexterity', label: 'Geschicklichkeit' },
            { key: 'stamina', label: 'Widerstand' },
            { key: 'charisma', label: 'Charisma' },
            { key: 'manipulation', label: 'Manipulation' },
            { key: 'composure', label: 'Coolness' },
            { key: 'intelligence', label: 'Intelligenz' },
            { key: 'wits', label: 'Schläue' },
            { key: 'resolve', label: 'Entschlossenheit' },
          ],
        },
      ],
    },
    {
      title: 'Zustand',
      icon: 'lucide:droplet',
      fields: [
        { key: 'health', label: 'Gesundheit', type: 'number', min: 0 },
        { key: 'healthMax', label: 'Max. Gesundheit', type: 'number', min: 0 },
        { key: 'willpower', label: 'Willenskraft', type: 'number', min: 0 },
        { key: 'humanity', label: 'Menschlichkeit', type: 'number', min: 0, max: 10 },
        { key: 'hunger', label: 'Hunger', type: 'number', min: 0, max: 5 },
        { key: 'bloodPotency', label: 'Blutkraft', type: 'number', min: 0, max: 10 },
      ],
    },
    {
      title: 'Notizen',
      icon: 'lucide:notebook-pen',
      fields: [{ key: 'notes', label: 'Notizen', type: 'textarea', placeholder: 'Bindungen, Berührungen, Makel …' }],
    },
  ],
  subtitle(d) {
    if (d.subtitle) return d.subtitle
    const clan = d.clan ? String(d.clan).charAt(0).toUpperCase() + String(d.clan).slice(1) : ''
    return [clan, d.generation ? `Gen ${d.generation}` : ''].filter(Boolean).join(' · ') || '—'
  },
  cardStats(d) {
    return [
      { label: 'Gen', value: String(n(d, 'generation') || '—') },
      { label: 'Wille', value: String(n(d, 'willpower')) },
      { label: 'Blut', value: String(n(d, 'bloodPotency')) },
    ]
  },
  hpPercent(d) {
    const max = n(d, 'healthMax')
    if (max <= 0) return 100
    return Math.max(0, Math.min(100, (n(d, 'health') / max) * 100))
  },
}
