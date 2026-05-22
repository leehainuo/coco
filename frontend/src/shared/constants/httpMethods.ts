/**
 * HTTP method constant definitions
 */

import type { HttpMethod } from '../types'

/**
 * All supported HTTP methods
 */
export const HTTP_METHODS: HttpMethod[] = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS', 'TRACE']

/**
 * Tailwind CSS color class names for HTTP methods (semantic colors)
 */
export const HTTP_METHOD_COLORS: Record<HttpMethod, string> = {
  GET: 'bg-green-background text-green-foreground',
  POST: 'bg-orange-background text-orange-foreground',
  PUT: 'bg-blue-background text-blue-foreground',
  PATCH: 'bg-pink-background text-pink-foreground',
  DELETE: 'bg-red-background text-red-foreground',
  HEAD: 'bg-blue-background text-blue-foreground',
  OPTIONS: 'bg-blue-background text-blue-foreground',
  TRACE: 'bg-blue-background text-blue-foreground',
}

/**
 * Get the color class name for an HTTP method
 */
export function getMethodColor(method: string): string {
  const upperMethod = method?.toUpperCase() as HttpMethod
  return HTTP_METHOD_COLORS[upperMethod] || HTTP_METHOD_COLORS.GET
}
