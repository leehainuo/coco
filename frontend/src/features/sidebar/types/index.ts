/**
 * Sidebar feature module type definitions
 */

import type { Tag, Endpoint } from '@/shared/types'

export type { Tag, Endpoint }

/**
 * Sidebar state
 */
export interface SidebarState {
  search: string
  openTags: Set<string>
  selected: Endpoint | null
}
