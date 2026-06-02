import type { Character, Item, Memory, SceneSummary, StatTile } from '~/types/entities'

export function useDemoCampaign() {
  const campaign = {
    name: 'Asche & Aether',
    ruleset: 'dnd5e_2024' as const,
    session: 'Sitzung 14 · D&D 5e',
  }

  const stats: StatTile[] = [
    { key: 'characters', icon: 'lucide:users', value: '28', delta: '+3', deltaTone: 'up', accent: 'primary' },
    { key: 'scenes', icon: 'lucide:castle', value: '9', delta: '+1', deltaTone: 'up', accent: 'secondary' },
    { key: 'items', icon: 'lucide:gem', value: '54', delta: '+6', deltaTone: 'up', accent: 'gold' },
    { key: 'memories', icon: 'lucide:sparkles', value: '37', delta: '3', deltaTone: 'flat', accent: 'magenta' },
  ]

  const characters: Character[] = [
    {
      id: 'lyra', name: 'Lyra Nachtschatten', subtitle: 'Halbelf · Hexenmeisterin', type: 'pc', status: 'alive',
      initial: 'L', ring: 'conic-gradient(from 140deg,var(--primary),var(--secondary),var(--magenta),var(--primary))',
      stats: [{ label: 'RK', value: '17' }, { label: 'TP', value: '42' }, { label: 'Stufe', value: '7' }], hpPercent: 86,
    },
    {
      id: 'thorgan', name: 'Thorgan Eisenfaust', subtitle: 'Zwerg · Kämpfer', type: 'pc', status: 'alive',
      initial: 'T', ring: 'conic-gradient(from 120deg,var(--primary),var(--emerald),var(--primary))',
      stats: [{ label: 'RK', value: '19' }, { label: 'TP', value: '68' }, { label: 'Stufe', value: '8' }], hpPercent: 100,
    },
    {
      id: 'kael', name: 'Kael der Verlorene', subtitle: 'Mensch · Paladin · gefallen', type: 'pc', status: 'dead',
      initial: 'K', ring: 'conic-gradient(from 140deg,var(--ember),var(--gold),var(--magenta),var(--ember))',
      stats: [{ label: 'RK', value: '—' }, { label: 'TP', value: '0' }, { label: 'Akt II', value: '†' }], hpPercent: 0,
    },
    {
      id: 'seraphine', name: 'Seraphine Vael', subtitle: 'Tiefling · Schurkin', type: 'pc', status: 'hunted',
      initial: 'S', ring: 'conic-gradient(from 90deg,var(--secondary),var(--magenta),var(--secondary))',
      stats: [{ label: 'RK', value: '15' }, { label: 'TP', value: '31' }, { label: 'Stufe', value: '6' }], hpPercent: 64,
    },
    {
      id: 'orin', name: 'Magister Orin', subtitle: 'Mensch · NSC · Unbekannt', type: 'npc', status: 'unknown',
      initial: 'M', ring: 'conic-gradient(from 200deg,var(--primary),var(--secondary),var(--primary))',
      stats: [{ label: 'RK', value: '13' }, { label: 'TP', value: '?' }, { label: 'Stufe', value: '—' }], hpPercent: 40,
    },
    {
      id: 'vesper', name: 'Vesper Crow', subtitle: 'Mensch · Vampir · Gejagt', type: 'npc', status: 'hunted',
      initial: 'V', ring: 'conic-gradient(from 60deg,var(--gold),var(--ember),var(--gold))',
      stats: [{ label: 'RK', value: '16' }, { label: 'TP', value: '55' }, { label: 'Blut', value: '4' }], hpPercent: 78, critical: true,
    },
    {
      id: 'mira', name: 'Mira Sonnenweber', subtitle: 'Mensch · Klerikerin', type: 'ally', status: 'alive',
      initial: 'M', ring: 'conic-gradient(from 30deg,var(--gold),var(--primary),var(--gold))',
      stats: [{ label: 'RK', value: '18' }, { label: 'TP', value: '49' }, { label: 'Stufe', value: '7' }], hpPercent: 92,
    },
    {
      id: 'grix', name: 'Grix Zahnfletscher', subtitle: 'Goblin · Späher', type: 'foe', status: 'alive',
      initial: 'G', ring: 'conic-gradient(from 220deg,var(--emerald),var(--gold),var(--emerald))',
      stats: [{ label: 'RK', value: '14' }, { label: 'TP', value: '21' }, { label: 'Stufe', value: '3' }], hpPercent: 70,
    },
    {
      id: 'nyx', name: 'Die Schwester Nyx', subtitle: 'Schatten · Endgegnerin', type: 'foe', status: 'unknown',
      initial: 'N', ring: 'conic-gradient(from 160deg,var(--secondary),var(--magenta),var(--void-2),var(--secondary))',
      stats: [{ label: 'RK', value: '20' }, { label: 'TP', value: '?' }, { label: 'Stufe', value: '—' }], hpPercent: 100,
    },
  ]

  const items: Item[] = [
    { id: 'klingenstab', name: 'Klingenstab des Aethers', type: 'Zauberstab · Waffe', rarity: 'epic', icon: 'lucide:wand-2', attuned: true, note: 'Reagiert auf Mondlicht.' },
    { id: 'mondkrone', name: 'Krone des Aschemonds', type: 'Wundersamer Gegenstand', rarity: 'legendary', icon: 'lucide:crown', attuned: true },
    { id: 'siegel', name: 'Siegel der Versunkenen', type: 'Schlüssel · Quest', rarity: 'artifact', icon: 'lucide:hexagon' },
    { id: 'trank', name: 'Trank der Schattenheilung', type: 'Trank', rarity: 'rare', icon: 'lucide:flask-conical' },
    { id: 'dolch', name: 'Flüsterdolch', type: 'Waffe · Finesse', rarity: 'uncommon', icon: 'lucide:sword' },
    { id: 'amulett', name: 'Amulett der Wacht', type: 'Halskette', rarity: 'rare', icon: 'lucide:gem', attuned: true },
    { id: 'fackel', name: 'Ewige Fackel', type: 'Ausrüstung', rarity: 'common', icon: 'lucide:flame' },
    { id: 'karte', name: 'Karte der tiefen Wege', type: 'Dokument', rarity: 'uncommon', icon: 'lucide:scroll' },
  ]

  const scene: SceneSummary = {
    id: 'krypta',
    title: 'Die Versunkene Krypta',
    description: 'Tropfendes Wasser, uraltes Arkanlicht und ein Siegel, das niemand brechen sollte.',
    act: 'Akt III',
    status: 'Vorbereitet',
    actors: ['L', 'T', 'S'],
    extra: 4,
    tone: 'arcane',
  }

  const scenes: SceneSummary[] = [
    scene,
    { id: 'taverne', title: 'Taverne „Zum Aschemond"', description: 'Sicherer Hafen der Gruppe — und ein Geheimgang im Keller.', act: 'Stammort', status: 'Aktiv', actors: ['M', 'T'], extra: 2, tone: 'ember' },
    { id: 'wald', title: 'Der Flüsterwald', description: 'Nebel zwischen uralten Stämmen; etwas beobachtet die Reisenden.', act: 'Akt II', status: 'Vorbereitet', actors: ['L', 'S', 'G'], extra: 1, tone: 'verdant' },
    { id: 'thron', title: 'Thronsaal der Schwester Nyx', description: 'Schwarzer Marmor, kaltes Licht — der letzte Akt erwartet euch.', act: 'Finale', status: 'Entwurf', actors: ['N'], extra: 0, tone: 'arcane' },
    { id: 'hafen', title: 'Der Sturmhafen', description: 'Salz, Seile und geschmuggelte Fracht unter Deck.', act: 'Akt I', status: 'Vorbereitet', actors: ['V', 'L'], extra: 3, tone: 'ember' },
  ]

  const memories: Memory[] = [
    {
      id: 'm1', title: 'Schwäche des Lich-Königs', body: 'Lyra erwähnte sie betrunken in Akt II — nicht vergessen!',
      level: 'warning', subjectType: 'character', subjectLabel: 'Lyra Nachtschatten', time: 'vor 2 Std',
    },
    {
      id: 'm2', title: 'Kael ist gefallen', body: 'Im Kampf gegen die Schattengarde gestorben. Status auf „Tot".',
      level: 'critical', subjectType: 'character', subjectLabel: 'Kael der Verlorene', time: 'gestern', acknowledged: true,
    },
    {
      id: 'm3', title: 'Klingenstab reagiert auf Mondlicht', body: 'Entfaltet bei Vollmond eine zweite Wirkung.',
      level: 'notice', subjectType: 'item', subjectLabel: 'Klingenstab des Aethers', time: 'vor 2 Tg',
    },
    {
      id: 'm4', title: 'Geheimgang in der Taverne', body: 'Hinter dem Fass im Keller — führt zur Krypta.',
      level: 'info', subjectType: 'scene', subjectLabel: 'Zum Aschemond', time: 'vor 3 Tg',
    },
  ]

  const criticalAlert = {
    title: 'Vesper Crow wird von der Inquisition gejagt',
    body: 'Bei Tageslicht in der Stadt droht sofortige Konfrontation — vor der nächsten Szene klären.',
    subjectLabel: 'Vesper Crow',
    moreCount: 2,
  }

  return { campaign, stats, characters, items, scene, scenes, memories, criticalAlert }
}
