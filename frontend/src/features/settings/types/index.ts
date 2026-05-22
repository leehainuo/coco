/**
 * Auth feature module type definitions
 */

/**
 * Authentication configuration
 */
export interface AuthConfig {
  type: 'basic' | 'bearer' | 'apiKey' | 'oauth2'
  username?: string
  password?: string
  token?: string
  apiKey?: string
  apiKeyLocation?: 'header' | 'query'
  apiKeyName?: string
}

/**
 * Authentication state
 */
export interface AuthState {
  isAuthenticated: boolean
  config: AuthConfig | null
}
