const KEYS = ['baseUrl', 'apiKey', 'campaignId', 'campaignName']

async function getConfig() {
  const c = await chrome.storage.sync.get(KEYS)
  return {
    baseUrl: (c.baseUrl || '').replace(/\/+$/, ''),
    apiKey: c.apiKey || '',
    campaignId: c.campaignId || '',
    campaignName: c.campaignName || '',
  }
}

async function apiFetch(path, opts = {}) {
  const { baseUrl, apiKey } = await getConfig()
  if (!baseUrl || !apiKey) return { ok: false, error: 'not_configured' }
  let res
  try {
    res = await fetch(baseUrl + path, {
      ...opts,
      headers: { 'X-API-Key': apiKey, ...(opts.headers || {}) },
    })
  } catch (e) {
    return { ok: false, error: 'network: ' + String((e && e.message) || e) }
  }
  let body = null
  try {
    body = await res.json()
  } catch {
    /* no body */
  }
  if (!res.ok) return { ok: false, error: (body && body.error) || 'HTTP ' + res.status, status: res.status }
  return { ok: true, data: body }
}

chrome.runtime.onMessage.addListener((msg, _sender, sendResponse) => {
  ;(async () => {
    if (msg.type === 'config') {
      sendResponse(await getConfig())
    } else if (msg.type === 'listCampaigns') {
      sendResponse(await apiFetch('/ext/campaigns'))
    } else if (msg.type === 'existsBatch') {
      const { campaignId } = await getConfig()
      if (!campaignId) {
        sendResponse({ ok: false, error: 'no_campaign' })
        return
      }
      sendResponse(
        await apiFetch(`/ext/campaigns/${campaignId}/assets/exists`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ refs: msg.refs || [] }),
        }),
      )
    } else if (msg.type === 'upload') {
      const { campaignId } = await getConfig()
      if (!campaignId) {
        sendResponse({ ok: false, error: 'no_campaign' })
        return
      }
      sendResponse(
        await apiFetch(`/ext/campaigns/${campaignId}/assets`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ ref: msg.ref, name: msg.name, mime: msg.mime, data: msg.data }),
        }),
      )
    } else {
      sendResponse({ ok: false, error: 'unknown_message' })
    }
  })()
  return true
})
