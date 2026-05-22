/**
 * API Detail feature module type definitions
 */

import type { Endpoint, Tag } from '@/shared/types'

export type { Endpoint, Tag }

/**
 * API detail state
 */
export interface ApiDetailState {
  endpoint: Endpoint | null
  debugOpen: boolean
}
