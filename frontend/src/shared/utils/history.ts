/**
 * History utility functions
 */

import type { HistoryItem } from '../types'

const MAX_ITEMS = 20

function key(id: string): string {
  return `coco:hist:${id}`
}

/**
 * History management
 */
export function createHistory(id: string) {
  function list(): HistoryItem[] {
    try {
      const raw = localStorage.getItem(key(id))
      return raw ? JSON.parse(raw) : []
    } catch {
      return []
    }
  }

  function push(item: Omit<HistoryItem, 'id' | 'timestamp'>): void {
    const items = list()
    items.unshift({ ...item, id: Date.now().toString(), timestamp: Date.now() })
    if (items.length > MAX_ITEMS) items.splice(MAX_ITEMS)
    try {
      localStorage.setItem(key(id), JSON.stringify(items))
    } catch {
      // Ignore storage errors
    }
  }

  function clear(): void {
    localStorage.removeItem(key(id))
  }

  return { list, push, clear }
}
