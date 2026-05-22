/**
 * OpenAPI Schema utility functions
 * Compatible with both OpenAPI 2.0 (Swagger) and OpenAPI 3.0
 */

import type { SchemaObject, OpenAPISpec } from '../types'

/**
 * Resolve schema type
 * Works with both OpenAPI 2.0 and 3.0 schemas
 */
export function resolveType(schema: SchemaObject | undefined): string {
  if (!schema) return 'any'
  if (schema.$ref) return schema.$ref.split('/').pop() || 'any'
  // OpenAPI 3.0 features
  if (schema.anyOf) return schema.anyOf.map(resolveType).join(' | ')
  if (schema.oneOf) return schema.oneOf.map(resolveType).join(' | ')
  if (schema.allOf) return schema.allOf.map(resolveType).join(' & ')
  if (schema.type === 'array') return `array[${resolveType(schema.items as SchemaObject)}]`
  const base = schema.type || (schema.properties ? 'object' : 'any')
  // OpenAPI 3.0 nullable, OpenAPI 2.0 uses x-nullable
  return schema.nullable || (schema as any)['x-nullable'] ? `${base} | null` : base
}

/**
 * Schema helper functions
 * Supports both OpenAPI 2.0 ($ref: "#/definitions/...") and 3.0 ($ref: "#/components/schemas/...")
 */
export function makeSchemaHelpers(spec: OpenAPISpec | null) {
  function resolveRef(schema: SchemaObject | undefined): SchemaObject {
    if (!schema) return {}
    if (schema.$ref && spec) {
      // Support both OpenAPI 2.0 and 3.0 reference formats
      // OpenAPI 2.0: #/definitions/ModelName
      // OpenAPI 3.0: #/components/schemas/ModelName
      const parts = schema.$ref.replace('#/', '').split('/')
      let n: any = spec
      for (const p of parts) {
        n = n?.[p]
      }
      return n ?? {}
    }
    return schema
  }

  function schemaToExample(schema: SchemaObject | undefined, depth = 0): any {
    if (!schema || depth > 5) return null
    const s = resolveRef(schema)
    if (s.example !== undefined) return s.example
    if (s.anyOf?.length) return schemaToExample(resolveRef(s.anyOf[0] as SchemaObject), depth + 1)
    if (s.oneOf?.length) return schemaToExample(resolveRef(s.oneOf[0] as SchemaObject), depth + 1)
    if (s.allOf?.length) {
      const m: Record<string, any> = {}
      for (const sub of s.allOf) {
        const r = resolveRef(sub as SchemaObject)
        const ex = schemaToExample(r, depth + 1)
        if (ex && typeof ex === 'object' && !Array.isArray(ex)) {
          Object.assign(m, ex)
        }
      }
      return Object.keys(m).length ? m : null
    }
    switch (s.type) {
      case 'object': {
        const obj: Record<string, any> = {}
        for (const [k, v] of Object.entries(s.properties ?? {})) {
          obj[k] = schemaToExample(v as SchemaObject, depth + 1)
        }
        return obj
      }
      case 'array':
        return [schemaToExample(s.items as SchemaObject, depth + 1)]
      case 'integer':
      case 'number':
        return 0
      case 'boolean':
        return false
      case 'string':
        return s.format === 'date-time' ? '2024-01-01T00:00:00Z' : ''
      default:
        return s.properties ? schemaToExample({ ...s, type: 'object' }, depth + 1) : null
    }
  }

  function isXmlSchema(schema: SchemaObject | undefined): boolean {
    if (!schema) return false
    const s = resolveRef(schema)
    // Check if schema has xml property or if properties have xml tags
    if (s.xml) return true
    if (s.properties) {
      return Object.values(s.properties).some((prop: any) => prop.xml)
    }
    return false
  }

  function schemaToXml(schema: SchemaObject | undefined, rootName = 'root', depth = 0): string {
    if (!schema || depth > 5) return ''
    const s = resolveRef(schema)
    const tagName = s.xml?.name || rootName

    if (s.type === 'object' && s.properties) {
      const children: string[] = []
      for (const [k, v] of Object.entries(s.properties)) {
        const propSchema = v as SchemaObject
        const propName = propSchema.xml?.name || k
        const example = propSchema.example

        if (example !== undefined) {
          children.push(`  <${propName}>${example}</${propName}>`)
        } else if (propSchema.type === 'string') {
          children.push(`  <${propName}></${propName}>`)
        } else if (propSchema.type === 'integer' || propSchema.type === 'number') {
          children.push(`  <${propName}>0</${propName}>`)
        } else if (propSchema.type === 'boolean') {
          children.push(`  <${propName}>false</${propName}>`)
        } else {
          children.push(`  <${propName}></${propName}>`)
        }
      }
      return `<${tagName}>\n${children.join('\n')}\n</${tagName}>`
    }

    return `<${tagName}>${s.example || ''}</${tagName}>`
  }

  function formatSchema(schema: SchemaObject | undefined, mediaType?: string): string {
    if (!schema) return ''
    // Check if this is XML based on media type or schema properties
    const isXml = mediaType?.includes('xml') || isXmlSchema(schema)
    if (isXml) {
      return schemaToXml(schema, 'XMLRequest')
    }
    return JSON.stringify(schemaToExample(resolveRef(schema)), null, 2)
  }

  return { resolveRef, schemaToExample, formatSchema, schemaToXml }
}
