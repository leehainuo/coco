/**
 * Debug Panel feature module type definitions
 */

import type { DebugRequest, DebugResponse, HistoryItem, FormRow } from '@/shared/types'

export type { DebugRequest, DebugResponse, HistoryItem, FormRow }

/**
 * Debug panel state
 */
export interface DebugPanelState {
  sending: boolean
  response: DebugResponse | null
  error: string | null
  history: HistoryItem[]
}

/**
 * Parameter row data
 */
export interface ParamRow {
  name: string
  value: string
  enabled: boolean
}
