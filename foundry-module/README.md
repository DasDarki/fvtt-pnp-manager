# Aetherwright — FoundryVTT-Modul

Verbindet sich (ausgehend) mit dem Aetherwright-Relay und führt Jobs im GM-Client aus:
Actors/Items/Szenen anlegen, Ordner lesen. FoundryVTT **v13**.

## Build (bun)

```bash
bun install
bun run build      # -> dist/ (module.json + aetherwright.js)
```

Der Ordner `dist/` ist das installierbare Modul. Für die lokale Entwicklung in Foundror
nach `Data/modules/aetherwright/` symlinken oder kopieren.

## Pairing

1. Im Aetherwright-Frontend: **Foundry-Anbindung → Token erzeugen** → `wsUrl` + Token.
2. In Foundry: **Spieleinstellungen → Modul-Einstellungen → Aetherwright** → Relay-URL und
   Pairing-Token eintragen. Das Modul verbindet sich automatisch (nur GM).

## Protokoll

WebSocket, JSON-Nachrichten `{ id, type, payload, error }`:

- Modul → Relay beim Connect: `{ type: "hello", payload: { world, version } }`
- Relay → Modul Jobs: `create_actor` · `discover_folders` · `ping`
- Modul → Relay Antwort: `{ id, type: "result", payload }` bzw. `{ id, type: "error", error }`

> MVP: `img` wird als externe URL gesetzt. Upload in die The-Forge-Assets via
> `FilePicker.upload` folgt (M4-Ausbau).
