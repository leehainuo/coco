/**
 * Request body type constant definitions
 */

import type { BodyType, SelectOption } from '../types'

/**
 * All supported request body types
 */
export const BODY_TYPES: SelectOption<BodyType>[] = [
  { value: 'none', label: 'None' },
  { value: 'json', label: 'JSON' },
  { value: 'xml', label: 'XML' },
  { value: 'form', label: 'x-www-form-urlencoded' },
  { value: 'formdata', label: 'form-data' },
  { value: 'text', label: 'Text' },
  { value: 'binary', label: 'Binary' },
  { value: 'graphql', label: 'GraphQL' },
]

/**
 * Text-based request body types (use text editor)
 */
export const TEXT_BODY_TYPES: BodyType[] = ['json', 'xml', 'text', 'binary', 'graphql']

/**
 * Form-based request body types (use form editor)
 */
export const FORM_BODY_TYPES: BodyType[] = ['form', 'formdata']

/**
 * Check if the type is text-based
 */
export function isTextType(type: BodyType): boolean {
  return TEXT_BODY_TYPES.includes(type)
}

/**
 * Check if the type is form-based
 */
export function isFormType(type: BodyType): boolean {
  return FORM_BODY_TYPES.includes(type)
}

/**
 * Get the display label for a body type
 */
export function getBodyTypeLabel(type: BodyType): string {
  return BODY_TYPES.find((t) => t.value === type)?.label ?? type
}
