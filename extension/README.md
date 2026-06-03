# Aetherwright Image Grabber (Chrome/Chromium Extension)

Fügt auf **chatgpt.com** an jedes generierte Bild einen Button hinzu, der das Bild mit
einem Klick in die Bibliothek deiner **aktiven Aetherwright-Kampagne** hochlädt — inkl.
Dedup (dasselbe Bild wird nicht doppelt gespeichert).

## Installation (entpackt)

1. In Aetherwright unter **Einstellungen → API-Keys** einen Key erstellen (Scope `images`).
   Der Key wird nur einmal angezeigt — kopieren.
2. Chrome → `chrome://extensions` → **Entwicklermodus** an → **Entpackte Erweiterung laden**
   → diesen `extension/`-Ordner wählen.
3. Auf das Extension-Icon klicken:
   - **API-Basis-URL**: dieselbe wie die Web-App-API, z. B. `https://pnp.dasdarki.de/api/v1`
     (lokal: `http://localhost:8080/v1`).
   - **API-Key**: den eben erstellten `awk_…`-Key.
   - **Verbinden & Kampagnen laden** → Domain-Berechtigung erlauben → **aktive Kampagne** wählen → **Speichern**.
4. Auf chatgpt.com über ein Bild hovern → Button **„Aetherwright"** oben links:
   - **Klick** → fragt per Dialog nach einem **Namen** (vorausgefüllt mit dem ChatGPT-Titel) und lädt hoch.
   - **Shift-Klick** → lädt sofort mit dem **ChatGPT-Namen** hoch (ohne Dialog).
   - Grüner Haken = in der Bibliothek (bzw. war bereits vorhanden). Umbenennen geht später im Bild-Atelier.

## Wie es funktioniert

- `content.js` findet Bilder via `img[src*="/backend-api/estuary/content"]`, liest die stabile
  `file_…`-ID aus der URL (Dedup-Schlüssel) und injiziert den Button. Es lädt das Bild
  **same-origin** auf chatgpt.com als Blob.
- `background.js` (Service Worker) schickt das Bild als Base64 mit dem Key-Header
  `X-API-Key` an `POST {base}/ext/campaigns/{id}/assets`. Da der Worker die Domain-
  Berechtigung hat, gibt es **kein CORS-Problem**.
- Der Backend-Endpunkt dedupliziert über die `file_…`-ID (`source_ref`).
- **„Bereits hochgeladen"-Status:** Beim Laden fragt die Extension die sichtbaren Bild-IDs
  gebündelt am Server ab (`POST …/assets/exists`) und markiert vorhandene Bilder sofort mit
  grünem Haken — geräteübergreifend, nicht nur aus dem lokalen Cache. Bei Kampagnenwechsel
  wird neu geprüft.

## Sicherheit

- Der API-Key darf ausschließlich **Bilder hochladen** und prüfen, ob ein Bild schon
  existiert — kein Zugriff auf andere Daten.
- Schlüssel werden im Backend nur als SHA-256-Hash gespeichert; widerrufbar unter
  **Einstellungen → API-Keys**.
