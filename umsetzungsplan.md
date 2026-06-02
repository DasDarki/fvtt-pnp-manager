# Aetherwright — Umsetzungsplan

> Konzeptioneller Gesamtplan: Architektur, Entitäten, Komponenten und Phasen-Roadmap.
> Ergänzt [concept.md](concept.md). Produktname **Aetherwright** (Teil der -wright-Suite).

---

## 1. Vision & Scope

Ein arkanes Kommando-Deck, das den GM-Workflow bündelt: Charaktere, Items, Szenen
und einfache Bilder **konzeptionieren → mit Dall·E bebildern → per Klick nach FoundryVTT
(The Forge) übertragen**, inklusive Notizen/Erinnerungen, globaler Suche und einer
auto-adaptiven Ordnerstruktur, die die Foundry-Ordner spiegelt.

**MVP-Annahme:** Single-User (du als Spielleiter), aber Multi-Device. Mehrere Kampagnen
parallel, je mit eigenem Regelwerk und eigenem Forge-Root-Ordner.

---

## 2. System-Architektur

```
┌───────────────────────────┐        ┌──────────────────────────────┐
│  Frontend (Nuxt 3 / Vue)  │  REST  │   Backend (Go / Fiber)       │
│  Tailwind · SCSS · Pinia  │◄──────►│   GORM · JWT · Postgres      │
│  vue-i18n · Command Deck  │  WS    │                              │
└───────────────────────────┘        │  ┌────────────────────────┐  │
            ▲                         │  │ Ruleset-Adapter-Layer  │  │
            │ (Live-Updates, opt.)    │  │ dnd5e_2024 · vampire_v5│  │
            ▼                         │  └────────────────────────┘  │
                                      │  ┌────────────────────────┐  │
   ┌──────────────────┐   Dall·E API  │  │ Dall·E-Service / Queue │  │
   │  OpenAI Images   │◄─────────────►│  └────────────────────────┘  │
   └──────────────────┘               │  ┌────────────────────────┐  │
                                      │  │ Foundry-Relay (WS-Hub) │  │
                                      │  └───────────▲────────────┘  │
                                      └──────────────┼───────────────┘
                                         outbound WS │ (Socket vom Modul)
                                      ┌──────────────┴───────────────┐
                                      │ FoundryVTT v13 @ The Forge   │
                                      │  Aetherwright-Modul (TS/Vite)│
                                      │  Actors · Items · Scenes ·   │
                                      │  Journals · Folders · Assets │
                                      └──────────────────────────────┘
```

**Kern-Entscheidung Foundry-Anbindung:** Das Foundry-Modul baut eine **ausgehende**
WebSocket-Verbindung zu unserem Relay auf (Pairing per Token). Damit brauchen wir keine
öffentliche Foundry-REST-API und kommen sauber durch The Forge. Das Backend schickt Jobs
(„erstelle Actor", „lade Bild hoch", „lies Ordnerbaum"), das Modul führt sie mit der
Foundry-API im GM-Client aus und meldet IDs/Pfade zurück.

### Repo-Layout (Monorepo)
```
/frontend          Nuxt 3
/backend           Go / Fiber
/foundry-module    TypeScript / Vite (eigenes Release-CI)
/deploy            docker-compose (Coolify)
/docs              dieses Dokument, Adapter-Specs
```

---

## 3. Entity-Modell

### Übersicht

| Entität            | Zweck                                              | Nach Foundry                         |
|--------------------|----------------------------------------------------|--------------------------------------|
| **Campaign**       | Wurzel: Regelwerk + Forge-Root-Ordner + Verbindung | (Kontext, kein eigenes Dokument)     |
| **Character**      | PC/NSC mit regelwerk-spezifischen Stats            | Actor                                |
| **Item**           | Gegenstand mit Seltenheit/Stats                    | Item                                 |
| **Scene**          | Szene/Map mit verknüpften Akteuren                 | Scene                                |
| **Image** („Bild") | Schlichtes Bild / leerer Actor                     | leerer Actor · Journal · Tile        |
| **Asset**          | Gespeicherte Bilddatei (Upload/Dall·E)             | Datei-Upload in Forge-Storage        |
| **Folder**         | Auto-adaptive Ordner, spiegelt Foundry             | Folder (je Typ: Actor/Item/Scene/…)  |
| **Memory** („Erinnerung") | Notiz mit Log-Level; *kritisch* wird prominent | optional Journal                  |
| **Tag**            | Querschnitts-Kategorisierung                       | —                                    |
| **DalleJob**       | Bild-Generierungsauftrag + Ergebnis                | —                                    |
| **User/Device/RefreshToken** | Auth, Multi-Device                       | —                                    |
| **FoundryConnection** | Pairing-Status des Moduls                        | —                                    |

### Gemeinsame Felder (Basis der „Inhalts"-Entitäten Character/Item/Scene/Image)
`id`, `campaign_id`, `folder_id`, `name`, `slug`, `summary`, `image_asset_id`,
`system_data` (JSONB, regelwerk-spezifisch), `foundry_uuid`, `foundry_doc_id`,
`sync_state` (none/pending/synced/dirty/error), `sort`, `created_at`, `updated_at`.

### Campaign (Wurzel)
```
id, owner_id
name, slug, description
ruleset          enum: dnd5e_2024 | vampire_v5   (erweiterbar)
forge_root_path  string   z.B. "worlds/<world>/aetherwright" / Asset-Root
cover_asset_id
connection_id    -> FoundryConnection
settings         jsonb   (Defaults: Token-Größen, Auto-Upload, etc.)
created_at, updated_at, archived_at
```

### Character
`+ character_type` (pc/npc/ally/foe/neutral) · `status` (alive/dead/unknown/hunted) ·
`system_data` (Abilities, Vitals, Features … pro Regelwerk) · `foundry_actor_id`.

### Item
`+ item_type` · `rarity` (common…artifact) · `attuned?` · `system_data` · `foundry_item_id`.

### Scene
`+ scene_status` (draft/prepared/active/archived) · `map_asset_id` · `grid` (jsonb) ·
`actors` (m:n über `scene_actor`) · `foundry_scene_id`.

### Image („Bild")
`+ asset_id` · `push_as` (empty_actor | journal | tile) · `notes`. Der leichteste Typ —
im Kern ein Asset mit Kontext, das als leerer Actor/Handout nach Foundry geht.

### Memory („Erinnerung") — polymorph, mit Log-Level
```
id, campaign_id
subject_type   character | item | scene | image | campaign
subject_id     (null = kampagnen-allgemein)
title, body (rich), kind (note/reminder/secret/event)
level          info | notice | warning | critical   (= Prio/Log-Level)
acknowledged   bool   (kritische quittierbar, ohne sie zu löschen)
pinned, created_at, updated_at
```

**`level` ist gleichzeitig Prio und Schweregrad** und steuert Sortierung **und** Sichtbarkeit:

| Level      | Farbe   | Anzeige                                                                 |
|------------|---------|------------------------------------------------------------------------|
| `info`     | dim     | normal im Feed                                                          |
| `notice`   | cyan    | normal, leicht hervorgehoben                                           |
| `warning`  | gold    | Badge auf Karte + im Feed oben                                          |
| `critical` | ember   | **überall prägnant**: Alarm-Strip im Dashboard, rote Eck-Ribbon auf der Entitäts-Karte, Marker in Suche/Command-Palette/Detail-Header |

**Kritische Hinweise (Architektur):**
- Dashboard zeigt oben einen **`CriticalAlerts`-Strip**, der alle nicht-quittierten
  `critical`-Memories der Kampagne aggregiert (genau gegen den Schmerz: „Charakter ist tot/Verräter"
  ohne Scrollen sofort sichtbar).
- Entitäts-Karten/Detail tragen ein `LevelBadge` bzw. eine glühende Eck-Ribbon bei `critical`.
- `acknowledged` blendet den Alarm aus, **ohne** die Notiz zu verlieren (Historie bleibt).
- Backend: `GET /campaigns/:id/alerts` liefert offene kritische Memories; optional Live-Push über `/ws/app`.

### Folder — auto-adaptiv
```
id, campaign_id, parent_id
name, color, sort
foundry_folder_id, foundry_type (Actor|Item|Scene|JournalEntry)
origin (foundry | local)   -> beim Push wird local → foundry gemappt
```

### Asset / DalleJob / Tag / Auth — siehe Schema in Phase M1/M3.

---

## 4. Regelwerk-Abstraktion (Adapter-Layer) — Architektur-Kern

Statt fester Felder pro Spiel definieren wir **System-Adapter**. Ein Adapter beschreibt
ein Regelwerk komplett und treibt sowohl die UI als auch den Foundry-Export:

```ts
interface RulesetAdapter {
  id: 'dnd5e_2024' | 'vampire_v5'
  label: I18nKey
  characterSchema: FieldGroup[]      // treibt das schema-getriebene Formular
  itemSchema: FieldGroup[]
  deriveStats(data): DerivedStats    // z.B. Ability-Mods, Rüstungsklasse
  toFoundryActor(character): FoundryActorDoc   // baut { type, system, prototypeToken }
  fromFoundryActor?(doc): Partial<Character>   // optionaler Import
}
```

- **Frontend:** Ein generischer `SchemaForm`-Renderer baut aus `characterSchema` die Eingabe-
  masken aus den Design-System-Primitives. Neues Regelwerk = neuer Adapter, **kein** neues UI.
- **Backend:** spiegelt Adapter für Validierung + `toFoundryActor`-Mapping
  (D&D 2024 → Foundry-System `dnd5e`; Vampire V5 → `vtm5e` o. ä.).
- `system_data` der Entitäten ist JSONB — exakt das Muster, das Foundry selbst nutzt.

**Erst-Adapter:** `dnd5e_2024` (D&D „5.5e"/2024) und `vampire_v5`.

---

## 5. Foundry-Sync (Modul + Relay)

1. **Pairing:** Web-App erzeugt Token → im Modul (Settings) eingetragen → Modul öffnet
   authentifizierte WS zum Relay. Status zurück an `FoundryConnection`.
2. **Ordner-Discovery (auto-adaptiv):** Modul liest `game.folders` (je Typ), schickt Baum →
   Backend mergt in `folders`. App nutzt/respektiert die Struktur.
3. **Push-Jobs:** `create_actor | create_item | create_scene | create_empty_actor |
   upload_asset | ensure_folder`. Modul legt im richtigen Ordner unterhalb `forge_root_path`
   an (fehlende Ordner werden erstellt) und gibt `uuid`/Pfad zurück → `foundry_*`-Felder.
4. **Bild-Upload:** Modul lädt Asset-Bytes über signierte URL vom Backend, lädt via Foundry
   `FilePicker.upload` in den Forge-Storage, setzt `actor.img`/Token-Textur.
5. **Reverse-Sync (optional, später):** Foundry-Hooks (`createFolder`, `updateActor`) → Backend.
6. **Release:** GitHub Actions baut das Modul, packt ZIP + `module.json` bei Semver-Tag.

> ⚠️ The-Forge-Spezifika (Asset-Library, Upload-Pfade) in M4 verifizieren.

---

## 6. Dall·E-Integration

- `DalleJob`: Prompt, Größe/Qualität/Style, Status, Ergebnis-Asset.
- UI: **Dall·E-Atelier** + Inline-Panel im Entity-Editor („Bild generieren").
- Flow: Prompt → Backend ruft OpenAI Images → Ergebnis als `Asset` (eigenes Storage zuerst)
  → an Entität hängen → bei Foundry-Push nach Forge spiegeln.
- Prompt-Vorlagen pro Entity-Typ (Portrait/Map/Handout), Verlauf, erneut generieren.

---

## 7. Frontend — Komponenten-Inventar

**Layout / Shell**
`AppShell` · `Sidebar` (+ `FolderTree`, `SyncStatus`, `CampaignSwitcher`, `UserMenu`) ·
`TopCommandBar` · `CommandPalette` (Ctrl+K) · `ThemeToggle` · `LangToggle`.

**Primitives (aus dem Design-System)**
`AwButton` · `AwInput` · `AwSelect` · `AwTextarea` · `AwSwitch` · `AwCheckbox` ·
`AwSearchBar` · `AwKbd` · `AwStatusPill` · `AwRarity` · `AwTag` · `AwToast` ·
`AwMeter` (HP/Mana) · `AwDropzone` · `AwPanel`.

**Entity-Komponenten**
`CharacterCard` · `ItemCard` · `SceneCard` · `ImageCard` · `MemoryCard` / `MemoryFeed` ·
`StatBlock` (+ `AbilityGrid`) · `EntityGrid` · `EntityDetail` · `LinkGraph` (Drag&Drop-Verknüpfen) ·
`CriticalAlerts` (Dashboard-Alarm-Strip) · `LevelBadge` · `LevelPicker` (Editor).

**Feature-Panels**
`SchemaForm` (adapter-getrieben) · `DalleStudio` · `FoundrySyncPanel` · `RarityPicker` ·
`FolderPicker` · `TagEditor`.

**Pinia-Stores**
`auth` · `campaign` · `entities` · `folders` · `memories` · `tags` · `search` · `dalle` ·
`foundry` (Sync) · `ui` (Theme/Lang/Palette/Hotkeys).

**Routen**
`/login` · `/campaigns` · `/c/:id/dashboard` · `/c/:id/characters[/:cid]` ·
`/c/:id/items` · `/c/:id/scenes` · `/c/:id/images` · `/c/:id/memories` ·
`/c/:id/dalle` · `/c/:id/settings/foundry`.

---

## 8. Backend — Struktur (Go / Fiber)

```
cmd/server
internal/
  config  httpx  auth(jwt+refresh, devices)
  campaign  entity  folder  memory  tag  asset  search
  ruleset/adapters   dalle   foundry(relay-ws, jobs, protocol)
  db (gorm models + migrations)
```
- REST unter `/api/v1`, WS `/ws/foundry` (Modul) und optional `/ws/app` (Live-Updates).
- JWT Access + Refresh, Multi-Device (RefreshToken je Device, Rotation).

---

## 9. Querschnitt

- **i18n:** DE Default, EN optional (vue-i18n). Google-Fonts im Build **lokal bundlen**
  (Datenschutz, [concept.md](concept.md)).
- **Security:** JWT-Rotation, scoped Asset-URLs (signiert), Pairing-Token mit Ablauf,
  Eingabe-Validierung serverseitig über Adapter.
- **UX:** Hotkeys/Command-Palette, Drag&Drop-Verknüpfen, Light/Dark (persistiert),
  responsive, `prefers-reduced-motion`.
- **Deployment:** docker-compose für Coolify; Modul-Release via GitHub-CI auf Semver-Tag.
- **Kein** Testing/Linting (laut concept.md).

---

## 10. Roadmap — so arbeiten wir es ab

> Empfehlung: **vertikaler Tracer-Bullet** — Campaign + Character einmal komplett bis nach
> Foundry, danach Item/Scene/Image nach gleichem Muster ergänzen.

**M0 · Fundament**
- [ ] Monorepo + docker-compose + Postgres + Coolify-Grundgerüst
- [ ] Design-Tokens → Tailwind-Config + SCSS aus `design-system.html`
- [ ] Nuxt-Setup, i18n (DE/EN), Fonts lokal gebundlet
- [ ] Auth (JWT + Refresh, Multi-Device), Login-Flow

**M1 · Kern-Daten & Dashboard (ohne Foundry)**
- [ ] DB-Schema + Migrationen (alle Entitäten, JSONB system_data)
- [ ] Campaign-CRUD (Regelwerk-Auswahl + Forge-Root-Ordner)
- [ ] Generisches Entity-CRUD: Character/Item/Scene/Image
- [ ] Folder-Baum (manuell), Tags, Memories (polymorph)
- [ ] Globale Suche
- [ ] Dashboard an echte Daten anschließen (aus `dashboard.html`)

**M2 · Regelwerk-Adapter & Editor**
- [ ] `RulesetAdapter`-Interface + `SchemaForm`-Renderer
- [ ] Adapter `dnd5e_2024` (Stats, StatBlock, deriveStats)
- [ ] Adapter `vampire_v5`
- [ ] Charakter-Detailansicht / Editor

**M3 · Dall·E**
- [ ] Dall·E-Service + Job-Queue + Asset-Storage
- [ ] Atelier-UI + Inline-Generierung im Editor + Prompt-Vorlagen

**M4 · Foundry-Modul & Sync**
- [ ] Modul-Skeleton (Vite, v13) + Pairing/Settings + WS-Relay
- [ ] Ordner-Discovery (auto-adaptiv) → App
- [ ] Push: Actor/Item/Scene/leerer Actor + Bild-Upload nach Forge
- [ ] `foundry_*`-Rücksync der IDs/Pfade
- [ ] GitHub-CI: ZIP + `module.json` auf Semver-Tag

**M5 · Politur & Release**
- [ ] Drag&Drop-Verknüpfen, Command-Palette-Aktionen, Hotkeys
- [ ] Responsive/A11y/Performance-Feinschliff
- [ ] Deployment Coolify, optional Reverse-Sync

---

## 11. Offene Entscheidungen

1. **D&D „5.5e"** = D&D-2024-Regeln → Foundry-System `dnd5e` (v4+)? Bestätigen.
2. **Vampire V5** in Foundry: Ziel-System/Actor-Typ (`vtm5e`?) — und wie verhält sich das
   zum erwähnten „Vicar"-Editor der anderen Gruppe?
3. **Single-User** für MVP bestätigt (keine Spieler-/Co-GM-Freigaben zunächst)?
4. **Asset-Storage:** eigenes Object-Storage (z.B. S3/MinIO) primär, dann Spiegelung nach
   Forge beim Push — ok?
