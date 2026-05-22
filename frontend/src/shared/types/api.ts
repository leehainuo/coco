/**
 * API related type definitions
 */

import type { ParameterObject, RequestBodyObject, ResponseObject } from './openapi'

/**
 * HTTP method type
 */
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'HEAD' | 'OPTIONS' | 'TRACE'

/**
 * Request body type
 */
export type BodyType = 'none' | 'json' | 'xml' | 'form' | 'formdata' | 'text' | 'binary' | 'graphql'

/**
 * Parameter location type
 */
export type ParamLocation = 'path' | 'query' | 'header' | 'cookie' | 'formData' | 'body'

/**
 * API endpoint interface
 */
export interface Endpoint {
  id: string
  path: string
  method: HttpMethod
  summary?: string
  description?: string
  tags?: string[]
  parameters?: ParameterObject[]
  requestBody?: RequestBodyObject
  responses?: Record<string, ResponseObject>
  deprecated?: boolean
  operationId?: string
}

/**
 * API tag group
 */
export interface Tag {
  name: string
  description?: string
  endpoints: Endpoint[]
}

/**
 * Request parameter interface
 */
export interface Parameter extends ParameterObject {
  enabled?: boolean
  value?: any
}

/**
 * Form row data
 */
export interface FormRow {
  name: string
  value: string
  type: 'text' | 'file'
  file?: File | null
  enabled: boolean
}

/**
 * Debug request interface
 */
export interface DebugRequest {
  method: HttpMethod
  url: string
  headers: Record<string, string>
  params: Record<string, string>
  body?: string | FormData
  bodyType: BodyType
}

/**
 * Debug response interface
 */
export interface DebugResponse {
  status: number
  statusText: string
  headers: Record<string, string>
  data: any
  time: number
  size?: number
}

/**
 * History item
 */
export interface HistoryItem {
  id: string
  timestamp: number
  request: DebugRequest
  response?: DebugResponse
  error?: string
}

/**
 * Select option interface
 */
export interface SelectOption<T = string> {
  value: T
  label: string
  disabled?: boolean
}
