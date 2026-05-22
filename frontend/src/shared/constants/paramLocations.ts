/**
 * Parameter location constant definitions
 */

import type { ParamLocation } from '../types'

/**
 * All supported parameter locations
 */
export const PARAM_LOCATIONS: ParamLocation[] = ['path', 'query', 'header', 'cookie', 'formData', 'body']

/**
 * Display label mapping for parameter locations
 */
export const PARAM_LOCATION_LABELS: Record<string, string> = {
  path: 'Path',
  query: 'Query',
  header: 'Header',
  cookie: 'Cookie',
  formData: 'Form Data',
  body: 'Body',
}

/**
 * Get the display label for a parameter location
 */
export function getParamLocationLabel(location: ParamLocation): string {
  return PARAM_LOCATION_LABELS[location] || location
}
