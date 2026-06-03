declare const game: any
declare const Hooks: any
declare const Actor: any
declare const Item: any
declare const Scene: any
declare const JournalEntry: any
declare const Folder: any
declare const ui: any
declare const window: any
declare const FilePicker: any
declare const ForgeVTT: any

const MODULE_ID = 'aetherwright'

interface AwMessage {
  id?: string
  type: string
  payload?: any
  error?: string
}

async function ingestImage(url: string): Promise<string> {
  if (!url || !/^https?:/i.test(url)) return url
  try {
    const res = await fetch(url)
    const blob = await res.blob()
    const ext = (blob.type.split('/')[1] || 'png').replace('svg+xml', 'svg')
    const file = new File([blob], `aw-${Date.now()}.${ext}`, { type: blob.type })
    const source = typeof ForgeVTT !== 'undefined' && ForgeVTT?.usingTheForge ? 'forgevtt' : 'data'
    const target = 'aetherwright'
    try {
      await FilePicker.createDirectory(source, target)
    } catch {
      /* directory exists */
    }
    const result = await FilePicker.upload(source, target, file, {})
    return result?.path || url
  } catch {
    return url
  }
}

let socket: WebSocket | null = null
let reconnectTimer: number | null = null
let manualClose = false

function getSettings() {
  return {
    wsUrl: String(game.settings.get(MODULE_ID, 'wsUrl') || ''),
    token: String(game.settings.get(MODULE_ID, 'token') || ''),
  }
}

function send(msg: AwMessage) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(msg))
  }
}

async function handleJob(msg: AwMessage) {
  try {
    switch (msg.type) {
      case 'ping':
        send({ id: msg.id, type: 'result', payload: { pong: true } })
        break

      case 'create_actor': {
        if (msg.payload?.img) msg.payload.img = await ingestImage(msg.payload.img)
        const actor = await Actor.create(msg.payload)
        send({ id: msg.id, type: 'result', payload: { uuid: actor.uuid, id: actor.id, img: actor.img } })
        break
      }

      case 'create_item': {
        if (msg.payload?.img) msg.payload.img = await ingestImage(msg.payload.img)
        const item = await Item.create(msg.payload)
        send({ id: msg.id, type: 'result', payload: { uuid: item.uuid, id: item.id, img: item.img } })
        break
      }

      case 'create_scene': {
        if (msg.payload?.background?.src) msg.payload.background.src = await ingestImage(msg.payload.background.src)
        const scene = await Scene.create(msg.payload)
        send({ id: msg.id, type: 'result', payload: { uuid: scene.uuid, id: scene.id } })
        break
      }

      case 'create_journal': {
        for (const page of msg.payload?.pages ?? []) {
          if (page.src) page.src = await ingestImage(page.src)
        }
        const journal = await JournalEntry.create(msg.payload)
        send({ id: msg.id, type: 'result', payload: { uuid: journal.uuid, id: journal.id } })
        break
      }

      case 'create_folder': {
        const folder = await Folder.create({
          name: msg.payload?.name,
          type: msg.payload?.type || 'Actor',
          color: msg.payload?.color || undefined,
          folder: msg.payload?.parent || null,
        })
        send({ id: msg.id, type: 'result', payload: { id: folder.id, uuid: folder.uuid } })
        break
      }

      case 'discover_folders': {
        const folders = (game.folders?.contents ?? []).map((f: any) => ({
          id: f.id,
          name: f.name,
          type: f.type,
          parent: f.folder?.id ?? null,
        }))
        send({ id: msg.id, type: 'result', payload: { folders } })
        break
      }

      default:
        send({ id: msg.id, type: 'error', error: `unknown job: ${msg.type}` })
    }
  } catch (e: any) {
    send({ id: msg.id, type: 'error', error: String(e?.message ?? e) })
  }
}

function connect() {
  const { wsUrl, token } = getSettings()
  if (!wsUrl || !token) {
    ui.notifications?.warn('Aetherwright: Relay-URL oder Pairing-Token fehlt.')
    return
  }
  manualClose = false
  try {
    socket?.close()
  } catch {
    /* ignore */
  }

  socket = new WebSocket(`${wsUrl}?token=${encodeURIComponent(token)}`)

  socket.addEventListener('open', () => {
    ui.notifications?.info('Aetherwright verbunden.')
    send({ type: 'hello', payload: { world: game.world?.id, version: game.version } })
  })

  socket.addEventListener('message', (ev: MessageEvent) => {
    let msg: AwMessage
    try {
      msg = JSON.parse(ev.data)
    } catch {
      return
    }
    handleJob(msg)
  })

  socket.addEventListener('close', () => {
    socket = null
    if (!manualClose) {
      if (reconnectTimer) window.clearTimeout(reconnectTimer)
      reconnectTimer = window.setTimeout(connect, 5000)
    }
  })

  socket.addEventListener('error', () => {
    try {
      socket?.close()
    } catch {
      /* ignore */
    }
  })
}

function disconnect() {
  manualClose = true
  if (reconnectTimer) window.clearTimeout(reconnectTimer)
  try {
    socket?.close()
  } catch {
    /* ignore */
  }
  socket = null
}

Hooks.once('init', () => {
  game.settings.register(MODULE_ID, 'wsUrl', {
    name: 'Relay-URL',
    hint: 'Die ws(s)://-Adresse aus dem Aetherwright-Atelier.',
    scope: 'world',
    config: true,
    type: String,
    default: '',
  })
  game.settings.register(MODULE_ID, 'token', {
    name: 'Pairing-Token',
    hint: 'Der einmalige Token aus den Foundry-Einstellungen deiner Kampagne.',
    scope: 'world',
    config: true,
    type: String,
    default: '',
    onChange: () => {
      if (game.user?.isGM) connect()
    },
  })
})

Hooks.once('ready', () => {
  if (game.user?.isGM) connect()
  ;(game.modules.get(MODULE_ID) as any).api = { connect, disconnect }
})
