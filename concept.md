## Problem & Kontext:
- Ich nutze für meine PnP Gruppe FoundryVTT als Plattform um zu spielen, epxlizit The Forge (hosting service)
- Ebenfalls genutz wird: DnD Beyond als Characterverwaltung, Obisidian für mich als Notizen und Storysammlung und ChatGPT bzw Dall-E für die Erstellung von bildbasierten Inhalten
- Wir haben zwar noch eine andere Gruppe die auf Vampire V5 spielt und dabei ebenfalls FoundryVTT nur mit einem custom Character Editor namens Vicar
- Leider ist es sehr anstrengend in Obsidian Characters zu konzeptioeren, ein Bild in ChatGPT zu generieren, es dann in Obisidian zu hinterlegen, es in The Forge hoch zu laden um es dann Ingame in FoundryVTT als Actor anzulegen und dann mit dem Bild zu verbinden. Das sind einfach unheimlich viele Schritte. Zumal Notizen, Anmerkungen und Sachen die ich mir merken muss über Characters und Szenen, einfach in einem großen Notiz Dokument das nur pro Session sortiert ist. Wenn also ein Character gestorben ist, muss ich immer sicher gehen durchs gucken in dieser Notiz Liste um zu wissen, dass er tot ist. Das ist einfach unpraktisch und fehleranfällig.

## Lösungsidee:
- Eine Webanwendung, die es ermöglicht Characktere, Szenen, Items und andere kleine Informationen zu erstellen, zu verwalten und zu organisieren.
- Die Webanwendung hat eine Dall-E API Integration um direkt Bilder zu generieren
- Mit einem einfachen Klick und einem FoundryVTT Modul können Actors, Items und einfache Bilder (leerer Actor) direkt in FoundryVTT übertragen werden
- Characters können dabei basierend auf dem Regelwerk angelegt werden und mit Stats und Features hinterlegt werden, die dann auch in FoundryVTT übernommen werden
- Informationen, Anmerkungen und "Erinnerungen" können direkt an dem Character, dem Item, der Szene oder dem "generellen" Info/Bild angehangen werden
- Eine Suche über alles ermöglicht es schnell Informationen zu finden, ohne durch lange Notizen scrollen zu müssen
- Eine Auto-Adaptive Ordner Struktur erkennt die Ordner die bereits in FoundryVTT sind und ermöglicht es diese auch in der Webanwendung zu nutzen. Die Anwendung respektiert dann ebenfalls die Ordnerstruktur bei Erstellung neuer Actors, Items und Bilder in FoundryVTT

## Design:
- Modernes futuristisches Design
- Gradients, Smoke, Neon und coole kleine Animationen
- Responsive für Desktop und Mobile
- Intuitive und einfache Bedienung
- Übersichtliche Darstellung von Characters, Items, Szenen und Informationen
- Ein Dashboard, das einen schnellen Überblick über alle wichtigen Informationen bietet (z.B. aktuelle Charaktere, Items, Szenen, Erinnerungen, etc.)
- Hotkeys für häufig genutzte Funktionen (z.B. neues Character erstellen, neues Item erstellen, etc.) (Quick Palette und Suche in einem)
- Eine Drag & Drop Funktion um Characters, Items und Szenen einfach zu organisieren und zu verbinden
- Light und Dark Mode

## Richtlinien:
- Fokus auf Benutzerfreundlichkeit und intuitive Bedienung
- Integration mit FoundryVTT und Dall-E API
- Sicherheit und Datenschutz der Benutzerdaten (Wenn Nutzung von Google Fonts, dann im Build Prozess downloaden und bundlen)
- Keine Kommentare schreiben
- Deutsche Sprache als Hauptsprache, aber auch Englisch als Option anbieten (i18n)

## Technologien:
- Frontend: Nuxt 3, Tailwind CSS, Vue 3, SCSS, Pinia, vue-i18n, Composition API (defineModel nicht emit verwenden)
- Backend: Go, Fiber, GORM, JWT (Multi Device Support) mit Refresh Tokens, Dall-E API Integration, FoundryVTT API Integration
- Datenbank: PostgreSQL
- Deployment: Docker Compose für Coolify
- Kein Testing, kein Linting
- FoundryVTT Module: TypeScript, Vite fürs Bundling, FoundryVTT v13, Github CI/CD für automatisches Packing on Semver Tag