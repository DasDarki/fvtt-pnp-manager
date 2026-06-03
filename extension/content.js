const IMG_SELECTOR = 'img[src*="/backend-api/estuary/content"]'
const REF_RE = /[?&]id=(file[_-][A-Za-z0-9]+)/

let cfg = { campaignId: '', campaignName: '', configured: false }
const doneRefs = new Set()
const pendingCheck = new Set()
let checkTimer = null

function refOf(src) {
  const m = (src || '').match(REF_RE)
  return m ? m[1] : null
}

function queueCheck(ref) {
  if (!cfg.configured || doneRefs.has(ref)) return
  pendingCheck.add(ref)
  if (checkTimer) clearTimeout(checkTimer)
  checkTimer = setTimeout(flushChecks, 400)
}

async function flushChecks() {
  checkTimer = null
  if (!cfg.configured || pendingCheck.size === 0) return
  const refs = [...pendingCheck]
  pendingCheck.clear()
  const res = await chrome.runtime.sendMessage({ type: 'existsBatch', refs })
  if (!res || !res.ok || !res.data) return
  const existing = res.data.existing || []
  if (existing.length) {
    await markManyDone(existing)
    refreshButtons()
  }
}

function recheckAll() {
  if (!cfg.configured) return
  document.querySelectorAll('.aw-grab-btn').forEach((btn) => {
    if (btn.dataset.state === 'idle') queueCheck(btn.dataset.awRef)
  })
}

async function loadConfig() {
  const c = await chrome.runtime.sendMessage({ type: 'config' })
  cfg = {
    campaignId: c.campaignId,
    campaignName: c.campaignName,
    configured: !!(c.baseUrl && c.apiKey && c.campaignId),
  }
  doneRefs.clear()
  if (cfg.campaignId) {
    const key = 'done_' + cfg.campaignId
    const stored = await chrome.storage.local.get(key)
    ;(stored[key] || []).forEach((r) => doneRefs.add(r))
  }
  refreshButtons()
  recheckAll()
}

async function markDone(ref) {
  return markManyDone([ref])
}

async function markManyDone(refs) {
  refs.forEach((r) => doneRefs.add(r))
  if (!cfg.campaignId) return
  const key = 'done_' + cfg.campaignId
  const stored = await chrome.storage.local.get(key)
  const arr = stored[key] || []
  let changed = false
  for (const r of refs) {
    if (!arr.includes(r)) {
      arr.push(r)
      changed = true
    }
  }
  if (changed) await chrome.storage.local.set({ [key]: arr })
}

function blobToDataURL(blob) {
  return new Promise((resolve, reject) => {
    const r = new FileReader()
    r.onload = () => resolve(r.result)
    r.onerror = reject
    r.readAsDataURL(blob)
  })
}

function setState(btn, state, title) {
  btn.dataset.state = state
  if (title !== undefined) btn.title = title
}

async function grab(btn, ref, src, name) {
  if (!cfg.configured) {
    setState(btn, 'error', 'Aetherwright: erst im Extension-Popup konfigurieren')
    return
  }
  if (btn.dataset.state === 'loading') return
  setState(btn, 'loading', 'Lädt hoch …')
  try {
    const resp = await fetch(src)
    if (!resp.ok) throw new Error('Bild laden: HTTP ' + resp.status)
    const blob = await resp.blob()
    const data = await blobToDataURL(blob)
    const res = await chrome.runtime.sendMessage({ type: 'upload', ref, name, mime: blob.type, data })
    if (res && res.ok) {
      await markDone(ref)
      setState(btn, 'done', res.data && res.data.deduped ? 'Bereits in der Bibliothek' : 'In Bibliothek hochgeladen ✓')
    } else {
      setState(btn, 'error', 'Fehler: ' + ((res && res.error) || 'unbekannt'))
    }
  } catch (e) {
    setState(btn, 'error', 'Fehler: ' + String((e && e.message) || e))
  }
}

function makeButton(ref, src, name) {
  const btn = document.createElement('button')
  btn.className = 'aw-grab-btn'
  btn.type = 'button'
  btn.dataset.awRef = ref
  btn.dataset.state = doneRefs.has(ref) ? 'done' : 'idle'
  btn.title = 'Klick: mit Namen hochladen · Shift-Klick: ChatGPT-Name'
  btn.innerHTML = '<span class="aw-ic"></span><span class="aw-lbl">Aetherwright</span>'
  btn.addEventListener('click', (e) => {
    e.preventDefault()
    e.stopPropagation()
    if (btn.dataset.state === 'done' || btn.dataset.state === 'loading') return
    let finalName = name
    if (!e.shiftKey) {
      const entered = window.prompt('Name in der Aetherwright-Bibliothek:', name)
      if (entered === null) return
      finalName = entered.trim() || name
    }
    grab(btn, ref, src, finalName)
  })
  return btn
}

function inject(img) {
  const src = img.currentSrc || img.src
  const ref = refOf(src)
  if (!ref) return
  const container = img.parentElement
  if (!container || container.querySelector(':scope > .aw-grab-btn')) return
  if (getComputedStyle(container).position === 'static') container.style.position = 'relative'
  const name = (img.getAttribute('alt') || '').replace(/^Generiertes Bild:\s*/i, '').trim() || 'ChatGPT Bild'
  container.appendChild(makeButton(ref, src, name))
  queueCheck(ref)
}

function scan(root) {
  if (root.querySelectorAll) root.querySelectorAll(IMG_SELECTOR).forEach(inject)
  if (root.matches && root.matches(IMG_SELECTOR)) inject(root)
}

function refreshButtons() {
  document.querySelectorAll('.aw-grab-btn').forEach((btn) => {
    if (btn.dataset.state === 'loading') return
    btn.dataset.state = doneRefs.has(btn.dataset.awRef) ? 'done' : 'idle'
  })
}

const observer = new MutationObserver((muts) => {
  for (const m of muts) {
    for (const node of m.addedNodes) {
      if (node.nodeType === 1) scan(node)
    }
  }
})

;(async function init() {
  await loadConfig()
  scan(document)
  observer.observe(document.documentElement, { childList: true, subtree: true })
  chrome.storage.onChanged.addListener((_changes, area) => {
    if (area === 'sync') loadConfig()
  })
})()
