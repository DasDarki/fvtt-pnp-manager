const $ = (id) => document.getElementById(id)
const KEYS = ['baseUrl', 'apiKey', 'campaignId', 'campaignName']

function normBase(v) {
  return (v || '').trim().replace(/\/+$/, '')
}
function originPattern(u) {
  try {
    return new URL(u).origin + '/*'
  } catch {
    return null
  }
}
function setStatus(msg, ok) {
  const s = $('status')
  s.textContent = msg
  s.className = 'status ' + (ok ? 'ok' : 'err')
}

async function load() {
  const c = await chrome.storage.sync.get(KEYS)
  $('baseUrl').value = c.baseUrl || ''
  $('apiKey').value = c.apiKey || ''
  if (c.campaignId) {
    const opt = document.createElement('option')
    opt.value = c.campaignId
    opt.textContent = c.campaignName || c.campaignId
    opt.selected = true
    $('campaign').appendChild(opt)
  }
}

async function ensurePermission(base) {
  const origin = originPattern(base)
  if (!origin) return false
  if (await chrome.permissions.contains({ origins: [origin] })) return true
  return chrome.permissions.request({ origins: [origin] })
}

async function connect() {
  const base = normBase($('baseUrl').value)
  const key = $('apiKey').value.trim()
  if (!base || !key) {
    setStatus('Basis-URL und Key eingeben.', false)
    return
  }
  await chrome.storage.sync.set({ baseUrl: base, apiKey: key })
  if (!(await ensurePermission(base))) {
    setStatus('Berechtigung für die Domain wurde nicht erteilt.', false)
    return
  }
  setStatus('Lade Kampagnen …', true)
  const res = await chrome.runtime.sendMessage({ type: 'listCampaigns' })
  if (!res || !res.ok) {
    setStatus('Fehler: ' + ((res && res.error) || 'unbekannt'), false)
    return
  }
  const sel = $('campaign')
  const cur = (await chrome.storage.sync.get('campaignId')).campaignId
  sel.innerHTML = '<option value="">—</option>'
  for (const cam of res.data) {
    const o = document.createElement('option')
    o.value = cam.id
    o.textContent = cam.name
    if (cam.id === cur) o.selected = true
    sel.appendChild(o)
  }
  setStatus(res.data.length + ' Kampagne(n) geladen.', true)
}

async function save() {
  const sel = $('campaign')
  const campaignId = sel.value
  const campaignName = (sel.options[sel.selectedIndex] && sel.options[sel.selectedIndex].textContent) || ''
  if (!campaignId) {
    setStatus('Bitte eine Kampagne wählen.', false)
    return
  }
  await chrome.storage.sync.set({
    baseUrl: normBase($('baseUrl').value),
    apiKey: $('apiKey').value.trim(),
    campaignId,
    campaignName,
  })
  setStatus('Gespeichert ✓ — aktiv: ' + campaignName, true)
}

$('connect').addEventListener('click', connect)
$('save').addEventListener('click', save)
load()
