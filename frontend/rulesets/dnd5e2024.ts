import type { RulesetAdapter, SystemData } from '~/types/ruleset'

function n(data: SystemData, key: string): number {
  const v = Number(data?.[key])
  return Number.isFinite(v) ? v : 0
}

export const dnd5e2024: RulesetAdapter = {
  id: 'dnd5e_2024',
  label: 'D&D 5e (2024)',
  schema: [
    {
      title: 'Grundlagen',
      icon: 'lucide:scroll-text',
      fields: [
        { key: 'race', label: 'Volk', type: 'text', placeholder: 'z.B. Halbelf' },
        { key: 'class', label: 'Klasse', type: 'text', placeholder: 'z.B. Hexenmeisterin' },
        { key: 'level', label: 'Stufe', type: 'number', min: 1, max: 20 },
        { key: 'background', label: 'Hintergrund', type: 'text', placeholder: 'z.B. Adlige' },
        {
          key: 'alignment',
          label: 'Gesinnung',
          type: 'select',
          options: [
            { value: '', label: '—' },
            { value: 'lg', label: 'Rechtschaffen Gut' },
            { value: 'ng', label: 'Neutral Gut' },
            { value: 'cg', label: 'Chaotisch Gut' },
            { value: 'ln', label: 'Rechtschaffen Neutral' },
            { value: 'tn', label: 'Neutral' },
            { value: 'cn', label: 'Chaotisch Neutral' },
            { value: 'le', label: 'Rechtschaffen Böse' },
            { value: 'ne', label: 'Neutral Böse' },
            { value: 'ce', label: 'Chaotisch Böse' },
          ],
        },
        { key: 'subtitle', label: 'Untertitel (Karte)', type: 'text', placeholder: 'Halbelf · Hexenmeisterin' },
      ],
    },
    {
      title: 'Attribute',
      icon: 'lucide:dices',
      fields: [
        {
          key: 'abilities',
          label: 'Attributswerte',
          type: 'abilities',
          traits: [
            { key: 'str', label: 'STR' },
            { key: 'dex', label: 'GES' },
            { key: 'con', label: 'KON' },
            { key: 'int', label: 'INT' },
            { key: 'wis', label: 'WIS' },
            { key: 'cha', label: 'CHA' },
          ],
        },
      ],
    },
    {
      title: 'Vitalwerte',
      icon: 'lucide:heart-pulse',
      fields: [
        { key: 'ac', label: 'Rüstungsklasse', type: 'number', min: 0 },
        { key: 'hp', label: 'Aktuelle TP', type: 'number', min: 0 },
        { key: 'hpMax', label: 'Maximale TP', type: 'number', min: 0 },
        { key: 'hpTemp', label: 'Temporäre TP', type: 'number', min: 0 },
        { key: 'speed', label: 'Bewegung (ft)', type: 'number', min: 0 },
        { key: 'initiative', label: 'Initiative-Bonus', type: 'number' },
        { key: 'proficiency', label: 'Übungsbonus', type: 'number', min: 0 },
        {
          key: 'spellcasting',
          label: 'Zauber-Attribut',
          type: 'select',
          options: [
            { value: '', label: '—' },
            { value: 'int', label: 'Intelligenz' },
            { value: 'wis', label: 'Weisheit' },
            { value: 'cha', label: 'Charisma' },
          ],
        },
      ],
    },
    {
      title: 'Rettungswürfe',
      icon: 'lucide:shield',
      fields: [
        {
          key: 'saves',
          label: 'Rettungswurf-Übungen',
          type: 'proficiencies',
          proficiencyMax: 1,
          traits: [
            { key: 'str', label: 'Stärke', ability: 'str' },
            { key: 'dex', label: 'Geschicklichkeit', ability: 'dex' },
            { key: 'con', label: 'Konstitution', ability: 'con' },
            { key: 'int', label: 'Intelligenz', ability: 'int' },
            { key: 'wis', label: 'Weisheit', ability: 'wis' },
            { key: 'cha', label: 'Charisma', ability: 'cha' },
          ],
        },
      ],
    },
    {
      title: 'Fertigkeiten',
      icon: 'lucide:target',
      fields: [
        {
          key: 'skills',
          label: 'Fertigkeiten (Klick: keine → geübt → Experte)',
          type: 'proficiencies',
          proficiencyMax: 2,
          traits: [
            { key: 'acr', label: 'Akrobatik', ability: 'dex' },
            { key: 'ani', label: 'Mit Tieren umgehen', ability: 'wis' },
            { key: 'arc', label: 'Arkane Kunde', ability: 'int' },
            { key: 'ath', label: 'Athletik', ability: 'str' },
            { key: 'dec', label: 'Täuschung', ability: 'cha' },
            { key: 'his', label: 'Geschichte', ability: 'int' },
            { key: 'ins', label: 'Einsicht', ability: 'wis' },
            { key: 'itm', label: 'Einschüchtern', ability: 'cha' },
            { key: 'inv', label: 'Nachforschung', ability: 'int' },
            { key: 'med', label: 'Medizin', ability: 'wis' },
            { key: 'nat', label: 'Naturkunde', ability: 'int' },
            { key: 'prc', label: 'Wahrnehmung', ability: 'wis' },
            { key: 'prf', label: 'Auftreten', ability: 'cha' },
            { key: 'per', label: 'Überzeugen', ability: 'cha' },
            { key: 'rel', label: 'Religion', ability: 'int' },
            { key: 'slt', label: 'Fingerfertigkeit', ability: 'dex' },
            { key: 'ste', label: 'Heimlichkeit', ability: 'dex' },
            { key: 'sur', label: 'Überleben', ability: 'wis' },
          ],
        },
      ],
    },
    {
      title: 'Währung',
      icon: 'lucide:coins',
      fields: [
        { key: 'pp', label: 'Platin (PP)', type: 'number', min: 0 },
        { key: 'gp', label: 'Gold (GP)', type: 'number', min: 0 },
        { key: 'ep', label: 'Elektrum (EP)', type: 'number', min: 0 },
        { key: 'sp', label: 'Silber (SP)', type: 'number', min: 0 },
        { key: 'cp', label: 'Kupfer (CP)', type: 'number', min: 0 },
      ],
    },
    {
      title: 'Notizen',
      icon: 'lucide:notebook-pen',
      fields: [{ key: 'notes', label: 'Biografie & Notizen', type: 'textarea', placeholder: 'Geheimnisse, Ziele, Bindungen …' }],
    },
  ],
  subtitle(d) {
    return d.subtitle || [d.race, d.class].filter(Boolean).join(' · ') || '—'
  },
  cardStats(d) {
    return [
      { label: 'RK', value: String(n(d, 'ac') || '—') },
      { label: 'TP', value: String(n(d, 'hp')) },
      { label: 'Stufe', value: String(n(d, 'level') || '—') },
    ]
  },
  hpPercent(d) {
    const max = n(d, 'hpMax')
    if (max <= 0) return 100
    return Math.max(0, Math.min(100, (n(d, 'hp') / max) * 100))
  },
}
