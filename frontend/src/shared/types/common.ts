/**
 * Common type definitions
 */

/**
 * Theme type
 */
export type Theme = 'light' | 'dark'

/**
 * Locale type
 */
export type Locale = 'zh' | 'en' | 'custom'

/**
 * Application configuration interface
 */
export interface AppConfig {
  theme: Theme
  locale: Locale
  sidebarWidth: number
  debugPanelWidth: number
}

/**
 * Size type
 */
export type Size = 'sm' | 'md' | 'lg'

/**
 * Variant type
 */
export type Variant = 'default' | 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive'

/**
 * Placement type
 */
export type Placement = 'top' | 'right' | 'bottom' | 'left' | 'top-start' | 'top-end' | 'bottom-start' | 'bottom-end'

/**
 * Loading state type
 */
export type LoadingState = 'idle' | 'loading' | 'success' | 'error'

/**
 * Nullable type
 */
export type Nullable<T> = T | null

/**
 * Optional type
 */
export type Optional<T> = T | undefined

/**
 * Deep readonly type
 */
export type DeepReadonly<T> = {
  readonly [P in keyof T]: T[P] extends object ? DeepReadonly<T[P]> : T[P]
}

/**
 * Deep partial type
 */
export type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P]
}
