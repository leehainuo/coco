/**
 * OpenAPI specification type definitions
 * Supports OpenAPI 3.0 and OpenAPI 2.0 (Swagger 2.0)
 */

export interface OpenAPISpec {
  openapi?: string
  swagger?: string
  info: InfoObject
  servers?: ServerObject[]
  paths: PathsObject
  components?: ComponentsObject
  basePath?: string
  host?: string
  schemes?: string[]
}

export interface InfoObject {
  title: string
  version: string
  description?: string
  termsOfService?: string
  contact?: ContactObject
  license?: LicenseObject
}

export interface ContactObject {
  name?: string
  url?: string
  email?: string
}

export interface LicenseObject {
  name: string
  url?: string
}

export interface ServerObject {
  url: string
  description?: string
  variables?: Record<string, ServerVariableObject>
}

export interface ServerVariableObject {
  enum?: string[]
  default: string
  description?: string
}

export type PathsObject = Record<string, PathItemObject>

export interface PathItemObject {
  $ref?: string
  summary?: string
  description?: string
  get?: OperationObject
  put?: OperationObject
  post?: OperationObject
  delete?: OperationObject
  options?: OperationObject
  head?: OperationObject
  patch?: OperationObject
  trace?: OperationObject
  servers?: ServerObject[]
  parameters?: (ParameterObject | ReferenceObject)[]
}

export interface OperationObject {
  tags?: string[]
  summary?: string
  description?: string
  externalDocs?: ExternalDocumentationObject
  operationId?: string
  parameters?: (ParameterObject | ReferenceObject)[]
  requestBody?: RequestBodyObject | ReferenceObject
  responses?: ResponsesObject
  callbacks?: Record<string, CallbackObject | ReferenceObject>
  deprecated?: boolean
  security?: SecurityRequirementObject[]
  servers?: ServerObject[]
}

export interface ExternalDocumentationObject {
  description?: string
  url: string
}

export interface ParameterObject {
  name: string
  in: 'query' | 'header' | 'path' | 'cookie'
  description?: string
  required?: boolean
  deprecated?: boolean
  allowEmptyValue?: boolean
  style?: string
  explode?: boolean
  allowReserved?: boolean
  schema?: SchemaObject | ReferenceObject
  example?: any
  examples?: Record<string, ExampleObject | ReferenceObject>
  content?: Record<string, MediaTypeObject>
  // Swagger 2.0
  type?: string
  format?: string
  items?: SchemaObject
  collectionFormat?: string
  default?: any
  maximum?: number
  exclusiveMaximum?: boolean
  minimum?: number
  exclusiveMinimum?: boolean
  maxLength?: number
  minLength?: number
  pattern?: string
  maxItems?: number
  minItems?: number
  uniqueItems?: boolean
  enum?: any[]
  multipleOf?: number
}

export interface RequestBodyObject {
  description?: string
  content: Record<string, MediaTypeObject>
  required?: boolean
}

export interface MediaTypeObject {
  schema?: SchemaObject | ReferenceObject
  example?: any
  examples?: Record<string, ExampleObject | ReferenceObject>
  encoding?: Record<string, EncodingObject>
}

export interface EncodingObject {
  contentType?: string
  headers?: Record<string, HeaderObject | ReferenceObject>
  style?: string
  explode?: boolean
  allowReserved?: boolean
}

export type ResponsesObject = Record<string, ResponseObject | ReferenceObject>

export interface ResponseObject {
  description: string
  headers?: Record<string, HeaderObject | ReferenceObject>
  content?: Record<string, MediaTypeObject>
  links?: Record<string, LinkObject | ReferenceObject>
}

export interface CallbackObject {
  [expression: string]: PathItemObject
}

export interface ExampleObject {
  summary?: string
  description?: string
  value?: any
  externalValue?: string
}

export interface LinkObject {
  operationRef?: string
  operationId?: string
  parameters?: Record<string, any>
  requestBody?: any
  description?: string
  server?: ServerObject
}

export interface HeaderObject extends Omit<ParameterObject, 'name' | 'in'> {}

export interface SchemaObject {
  title?: string
  multipleOf?: number
  maximum?: number
  exclusiveMaximum?: boolean
  minimum?: number
  exclusiveMinimum?: boolean
  maxLength?: number
  minLength?: number
  pattern?: string
  maxItems?: number
  minItems?: number
  uniqueItems?: boolean
  maxProperties?: number
  minProperties?: number
  required?: string[]
  enum?: any[]
  type?: string
  allOf?: (SchemaObject | ReferenceObject)[]
  oneOf?: (SchemaObject | ReferenceObject)[]
  anyOf?: (SchemaObject | ReferenceObject)[]
  not?: SchemaObject | ReferenceObject
  items?: SchemaObject | ReferenceObject
  properties?: Record<string, SchemaObject | ReferenceObject>
  additionalProperties?: boolean | SchemaObject | ReferenceObject
  description?: string
  format?: string
  default?: any
  nullable?: boolean
  discriminator?: DiscriminatorObject
  readOnly?: boolean
  writeOnly?: boolean
  xml?: XMLObject
  externalDocs?: ExternalDocumentationObject
  example?: any
  deprecated?: boolean
  $ref?: string
}

export interface DiscriminatorObject {
  propertyName: string
  mapping?: Record<string, string>
}

export interface XMLObject {
  name?: string
  namespace?: string
  prefix?: string
  attribute?: boolean
  wrapped?: boolean
}

export interface ReferenceObject {
  $ref: string
}

export interface ComponentsObject {
  schemas?: Record<string, SchemaObject | ReferenceObject>
  responses?: Record<string, ResponseObject | ReferenceObject>
  parameters?: Record<string, ParameterObject | ReferenceObject>
  examples?: Record<string, ExampleObject | ReferenceObject>
  requestBodies?: Record<string, RequestBodyObject | ReferenceObject>
  headers?: Record<string, HeaderObject | ReferenceObject>
  securitySchemes?: Record<string, SecuritySchemeObject | ReferenceObject>
  links?: Record<string, LinkObject | ReferenceObject>
  callbacks?: Record<string, CallbackObject | ReferenceObject>
}

export interface SecuritySchemeObject {
  type: 'apiKey' | 'http' | 'oauth2' | 'openIdConnect'
  description?: string
  name?: string
  in?: 'query' | 'header' | 'cookie'
  scheme?: string
  bearerFormat?: string
  flows?: OAuthFlowsObject
  openIdConnectUrl?: string
}

export interface OAuthFlowsObject {
  implicit?: OAuthFlowObject
  password?: OAuthFlowObject
  clientCredentials?: OAuthFlowObject
  authorizationCode?: OAuthFlowObject
}

export interface OAuthFlowObject {
  authorizationUrl?: string
  tokenUrl?: string
  refreshUrl?: string
  scopes: Record<string, string>
}

export type SecurityRequirementObject = Record<string, string[]>
