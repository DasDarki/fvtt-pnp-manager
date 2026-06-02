# Aetherwright

Arkanes Kommando-Deck für PnP-Kampagnen: Charaktere, Items, Szenen und Erinnerungen
konzeptionieren, mit Dall·E bebildern und per Klick nach **FoundryVTT** übertragen.
Teil der -wright-Suite. Siehe [Konzept](concept.md) und [Umsetzungsplan](umsetzungsplan.md).

## Struktur

```
frontend/        Nuxt 3 · Tailwind · SCSS · Pinia · vue-i18n  (Port 3000)
backend/         Go · Fiber · GORM · PostgreSQL · JWT          (Port 8080)
foundry-module/  FoundryVTT-Modul (TS/Vite, v13) — Relay-Client
deploy via       docker-compose.yml (Root) — db + backend + frontend
```

## Funktionen

- **Auth** (JWT Access + Refresh, Multi-Device)
- **Entitäten** Charakter · Item · Szene — je mit Editor, Dall·E-Bild und Erinnerungen
- **Regelwerk-Adapter** (D&D 5e 2024 an DDB/Foundry-`dnd5e` ausgerichtet, Vampire V5)
- **Bild-Atelier** mit **globalem Kampagnen-Stil** + **BYOK pro Benutzer** (OpenAI · Stability · Replicate · Imagen; Mock-Fallback ohne Key)
- **Erinnerungen** mit Log-Level (info/notice/warning/**critical**) + Dashboard-Alarm
- **FoundryVTT-Sync** über ausgehendes WS-Relay (Actor/Item/Scene + Bild-Upload in Forge-Assets)
- **Command-Palette** (Ctrl+K) — Suche über alles + Schnellaktionen · Light/Dark · DE/EN

## Lokale Entwicklung

```bash
# Backend + DB
cd backend && docker compose up -d db && go run ./cmd/server

# Frontend
cd frontend && bun install && bun run dev   # http://localhost:3000

# Foundry-Modul (Build -> dist/ ist installierbar)
cd foundry-module && bun install && bun run build
```

## Deployment (Docker Compose / Coolify)

Alles in Containern — so läuft es auch auf Coolify:

```bash
cp .env.example .env     # Secrets + öffentliche URLs setzen
docker compose up --build
```

- **Frontend** → `http://localhost:3000`, **Backend** → `http://localhost:8080`
- **Bild-KI (BYOK):** jeder Nutzer hinterlegt eigene Keys in der App unter **Einstellungen → Bild-Provider** (verschlüsselt in der DB, nie an den Browser zurückgegeben). `OPENAI_API_KEY` ist nur ein optionaler globaler Fallback (leer = Mock-Platzhalter).

### Coolify (Pfad-Routing auf einer Domain)

Setup: `pnp.dasdarki.de` → Frontend, `pnp.dasdarki.de/api` → Backend (Traefik strippt `/api`).

1. Repo als **Docker-Compose**-Ressource hinzufügen (`docker-compose.yml` im Root).
2. Domains in Coolify zuordnen:
   - Service **frontend** → `https://pnp.dasdarki.de` (Port 3000)
   - Service **backend** → `https://pnp.dasdarki.de/api` (Port 8080) — Coolify fügt für den
     Pfad automatisch eine **StripPrefix**-Middleware hinzu (`/api` wird entfernt).
3. Environment-Variablen setzen:
   - `PUBLIC_APP_URL=https://pnp.dasdarki.de`
   - `PUBLIC_API_URL=https://pnp.dasdarki.de/api`  ← **inkl. `/api`-Präfix**
   - `JWT_SECRET`, `ENCRYPTION_KEY` (lange Zufallswerte), `POSTGRES_PASSWORD`.
   Daraus ergeben sich automatisch: `NUXT_PUBLIC_API_BASE=…/api/v1`,
   `PUBLIC_BASE_URL=…/api` (Uploads & WS-Relay), `CORS_ORIGINS=PUBLIC_APP_URL`.
4. `pgdata` und `uploads` sind benannte Volumes (persistente DB + generierte Bilder).

> **Warum das ohne 404 funktioniert:** Das Backend lauscht auf **beiden** Präfixen
> `/api/v1/*` **und** `/v1/*`. Nach dem Strip kommt `/api/v1/…` als `/v1/…` an → trifft.
> Uploads (`/uploads`) und WS (`/ws/foundry`) liegen auf Root und werden über
> `…/api/uploads` bzw. `wss://…/api/ws/foundry` aufgerufen → nach Strip ebenfalls korrekt.
> Container-Healthcheck: `GET /health` (Root), `/api/v1/health` und `/v1/health`.

> Hinweis: `NUXT_PUBLIC_API_BASE` und `PUBLIC_BASE_URL` müssen **öffentliche** URLs sein
> (der Browser ruft API und Bilder direkt auf), nicht der interne `db`/`backend`-Hostname.

> Alternative ohne Strip: getrennte Subdomain `PUBLIC_API_URL=https://api.pnp.dasdarki.de`
> (Backend ohne Pfad-Präfix) — dann greift `/api/v1` direkt, ebenfalls unterstützt.
