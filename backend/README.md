# Aetherwright — Backend

Go · Fiber · GORM · PostgreSQL · JWT (Access + Refresh, Multi-Device). Teil der
[-wright-Suite](../umsetzungsplan.md).

## Entwicklung

Postgres per Docker, API lokal:

```bash
docker compose up -d db
cp .env.example .env
go run ./cmd/server
```

Oder alles in Containern (so läuft es auch auf Coolify):

```bash
docker compose up --build
```

API: `http://localhost:8080/api/v1`

## Endpunkte (M1)

| Methode | Pfad | Schutz | Zweck |
|--------|------|--------|-------|
| GET    | `/health` | — | Healthcheck |
| POST   | `/auth/register` | — | Registrierung → Token-Paar |
| POST   | `/auth/login` | — | Login → Token-Paar |
| POST   | `/auth/refresh` | — | Refresh-Token rotieren |
| POST   | `/auth/logout` | — | Refresh-Token widerrufen |
| GET    | `/auth/me` | Bearer | Aktueller Nutzer |
| GET/POST | `/campaigns` | Bearer | Kampagnen auflisten / anlegen |
| GET/PATCH/DELETE | `/campaigns/:id` | Bearer | Kampagne lesen / ändern / löschen |
| GET/POST | `/campaigns/:campaignId/characters` | Bearer | Charaktere |
| GET/PATCH/DELETE | `/campaigns/:campaignId/characters/:id` | Bearer | Charakter-Detail |
| GET/POST | `/campaigns/:campaignId/memories` | Bearer | Erinnerungen (mit `level`) |
| PATCH/DELETE | `/campaigns/:campaignId/memories/:id` | Bearer | u.a. `acknowledged` setzen |
| GET    | `/campaigns/:campaignId/alerts` | Bearer | Offene **kritische** Erinnerungen |

## Auth-Modell

- **Access-Token**: kurzlebiges JWT (HS256, Default 15 min), `Authorization: Bearer <token>`.
- **Refresh-Token**: zufälliger 256-bit-Token, serverseitig **nur als SHA-256-Hash**
  gespeichert (`refresh_tokens`), pro Gerät (`device`), Default 30 Tage. Bei `refresh`
  wird der alte Token widerrufen und ein neuer ausgegeben (Rotation) → Multi-Device.

## Konventionen

- UUID-Primärschlüssel (in Go generiert, kein DB-Extension-Zwang).
- `system_data`/`grid`/`settings` als JSONB (regelwerk-spezifisch).
- Alles ownership-scoped: Zugriff nur auf eigene Kampagnen und deren Inhalte.
- `Memory.level` (info/notice/warning/critical) treibt Prio **und** die prominente
  Anzeige im Frontend; `/alerts` liefert die offenen kritischen.
