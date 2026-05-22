/**
 * Settings composable
 */

import { ref, watch } from 'vue'

const STORAGE_KEY = 'coco:settings'

interface AuthData {
  type: 'none' | 'bearer' | 'apikey' | 'basic'
  bearer: {
    token: string
  }
  apikey: {
    name: string
    value: string
    in: 'header' | 'query'
  }
  basic: {
    username: string
    password: string
  }
}

const def: AuthData = {
  type: 'none',
  bearer: { token: '' },
  apikey: { name: 'X-API-Key', value: '', in: 'header' },
  basic: { username: '', password: '' },
}

function loadSettings(): AuthData {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) {
      return { ...def, ...JSON.parse(raw) }
    }
  } catch {
    // Ignore parse errors
  }
  return { ...def }
}

const settings = ref<AuthData>(loadSettings())

watch(
  settings,
  (value) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(value))
  },
  { deep: true }
)

export function useSettings() {
  /**
   * Apply settings to request headers and query parameters
   */
  function applySettings(headers: Record<string, string>, searchParams: URLSearchParams): void {
    const a = settings.value

    if (a.type === 'bearer' && a.bearer.token) {
      headers['Authorization'] = `Bearer ${a.bearer.token}`
    } else if (a.type === 'apikey' && a.apikey.value) {
      if (a.apikey.in === 'query') {
        searchParams.set(a.apikey.name || 'X-API-Key', a.apikey.value)
      } else {
        headers[a.apikey.name || 'X-API-Key'] = a.apikey.value
      }
    } else if (a.type === 'basic' && (a.basic.username || a.basic.password)) {
      headers['Authorization'] = `Basic ${btoa(`${a.basic.username}:${a.basic.password}`)}`
    }
  }

  return {
    settings,
    applySettings,
  }
}
