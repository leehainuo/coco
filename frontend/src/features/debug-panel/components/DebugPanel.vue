<template>
  <div
    :style="mobile ? {} : { width: width + 'px' }"
    :class="[
      'relative flex flex-col overflow-hidden',
      mobile
        ? 'border-border bg-background h-full w-full rounded-t-2xl border-x border-t'
        : 'border-border bg-background m-3 shrink-0 rounded-xl border shadow-2xl shadow-black/8 dark:shadow-black/40',
    ]"
  >
    <!-- Drag handle (mobile only) -->
    <div
      v-if="mobile"
      class="flex shrink-0 cursor-row-resize touch-none justify-center pt-2.5 pb-1.5"
      @mousedown.prevent="$emit('drag-start', $event)"
      @touchstart.prevent="$emit('drag-start', $event)"
    >
      <div class="bg-border h-1 w-10 rounded-full" />
    </div>
    <!-- Resize handle (desktop only) -->
    <div v-if="!mobile" @mousedown.prevent="startResize" class="absolute inset-y-0 left-0 z-10 w-1 cursor-col-resize" />

    <!-- Header -->
    <div class="border-border flex shrink-0 items-center justify-between border-b px-4 py-3">
      <span class="text-primary-foreground text-sm font-semibold">{{ t('onlineDebug') }}</span>
      <button
        @click="$emit('close')"
        class="text-icon-foreground hover:bg-icon-accent hover:text-icon-accent-foreground rounded p-1 transition-colors"
      >
        <XIcon :size="16" />
      </button>
    </div>

    <!-- URL bar -->
    <div class="border-border flex shrink-0 items-start gap-2 border-b px-3 py-2.5">
      <div class="flex min-w-0 flex-1 items-start gap-1.5">
        <span :class="['shrink-0 rounded-md px-1.5 py-1 text-xs font-bold uppercase', getMethodColor(endpoint.method)]">
          {{ endpoint.method }}
        </span>
        <p class="text-primary-foreground min-w-0 flex-1 truncate pt-1 font-mono text-xs">
          {{ url }}
        </p>
      </div>
      <Button variant="primary" size="sm" :loading="sending" @click="sendRequest">
        {{ sending ? t('sending') : t('send') }}
      </Button>
    </div>

    <!-- Scrollable params area -->
    <div class="scrollbar-hidden flex min-h-0 flex-1 flex-col overflow-y-auto py-4">
      <div class="divide-border divide-y">
        <!-- Path params -->
        <ParamSection v-if="path.length" :title="t('pathParams')" :items="path" @remove="(i) => path.splice(i, 1)" />

        <!-- Query params (default for GET / HEAD / DELETE / OPTIONS) -->
        <ParamSection
          v-if="isQueryVisible"
          :title="t('queryParams')"
          :items="query"
          :addable="true"
          @add="addQueryRow"
          @remove="(i) => query.splice(i, 1)"
        />

        <!-- Custom header -->
        <ParamSection
          v-if="isHeadersVisible"
          :title="t('header')"
          :items="headers"
          :addable="true"
          :name-editable="true"
          @add="headers.push({ name: '', value: '', enabled: true })"
          @remove="(i) => headers.splice(i, 1)"
        />

        <!-- Custom cookies -->
        <ParamSection
          v-if="isCookiesVisible"
          :title="t('cookie')"
          :items="cookies"
          :addable="true"
          :name-editable="true"
          @add="cookies.push({ name: '', value: '', enabled: true })"
          @remove="(i) => cookies.splice(i, 1)"
        />

        <!-- Body editor (default for POST / PUT / PATCH) -->
        <BodyEditor v-if="hasBody || isBodyVisible" ref="bodyEditorRef" v-model="body" v-model:bodyType="type" />

        <!-- History list -->
        <div v-if="isHistoryVisible" class="divide-border divide-y">
          <div class="flex shrink-0 items-center justify-between px-3 py-2">
            <span class="text-primary-foreground text-xs font-medium">{{ t('history') }}</span>
            <div class="flex items-center gap-2">
              <button
                @click="onClearHistory"
                class="text-muted-foreground hover:text-destructive text-[11px] transition-colors"
              >
                {{ t('historyClear') }}
              </button>
              <button
                @click="isHistoryVisible = false"
                class="text-icon-foreground hover:text-icon-accent-foreground transition-colors"
              >
                <XIcon :size="12" />
              </button>
            </div>
          </div>
          <p v-if="historyItems.length === 0" class="text-muted-foreground px-3 py-4 text-center text-xs">
            {{ t('historyEmpty') }}
          </p>
          <button
            v-for="item in historyItems"
            :key="item.id"
            class="hover:bg-accent flex w-full items-center gap-2 px-3 py-2 text-left transition-colors"
            @click="restoreHistory(item)"
          >
            <span
              :class="[
                'shrink-0 rounded px-1.5 py-0.5 text-[10px] font-bold',
                item.response?.status >= 200 && item.response?.status < 300
                  ? 'bg-green-background text-green-foreground'
                  : item.response?.status >= 400
                    ? 'bg-red-background text-red-foreground'
                    : 'bg-muted text-muted-foreground',
              ]"
            >
              {{ item.response?.status ?? '—' }}
            </span>
            <span class="text-muted-foreground shrink-0 font-mono text-xs">{{ fmtTime(item.ts) }}</span>
            <span class="text-muted-foreground ml-auto truncate text-[11px]">
              {{ item.response?.time }}ms · {{ item.response?.size }}
            </span>
          </button>
        </div>
      </div>

      <!-- Add section row (bottom, no border, hover bg per button) -->
      <div class="flex shrink-0 flex-wrap items-center gap-0.5 px-2 py-1.5 font-semibold">
        <button
          v-if="!isQueryVisible"
          @click="isQueryVisible = true"
          class="text-muted-foreground hover:bg-accent hover:text-accent-foreground inline-flex items-center gap-1 rounded px-2 py-1 text-xs transition-colors"
        >
          <PlusIcon :size="12" class="shrink-0" />
          {{ t('query') }}
        </button>
        <button
          v-if="!isHeadersVisible"
          @click="isHeadersVisible = true"
          class="text-muted-foreground hover:bg-accent hover:text-accent-foreground inline-flex items-center gap-1 rounded px-2 py-1 text-xs transition-colors"
        >
          <PlusIcon :size="12" class="shrink-0" />
          {{ t('header') }}
        </button>
        <button
          v-if="!isCookiesVisible"
          @click="isCookiesVisible = true"
          class="text-muted-foreground hover:bg-accent hover:text-accent-foreground inline-flex items-center gap-1 rounded px-2 py-1 text-xs transition-colors"
        >
          <PlusIcon :size="12" class="shrink-0" />
          {{ t('cookie') }}
        </button>
        <button
          v-if="!hasBody && !isBodyVisible"
          @click="isBodyVisible = true"
          class="text-muted-foreground hover:bg-accent hover:text-accent-foreground inline-flex items-center gap-1 rounded px-2 py-1 text-xs transition-colors"
        >
          <PlusIcon :size="12" class="shrink-0" />
          {{ t('body') }}
        </button>
        <button
          v-if="props.enableHistory && !isHistoryVisible"
          @click="onOpenHistory"
          class="text-muted-foreground hover:bg-accent hover:text-accent-foreground rounded px-2 py-1 text-xs transition-colors"
        >
          {{ t('history') }}
          <span v-if="historyItems.length" class="ml-1 opacity-50">{{ historyItems.length }}</span>
        </button>
      </div>
    </div>

    <!-- Divider / drag handle -->
    <div
      class="border-border hover:bg-accent flex h-2 shrink-0 cursor-row-resize items-center justify-center border-t transition-colors select-none"
      @mousedown.prevent="startResponseResize"
    >
      <span class="text-primary-foreground text-xs tracking-widest">···</span>
    </div>

    <!-- Response (height controlled by drag) -->
    <div :style="{ height: height + 'px' }" class="shrink-0 overflow-hidden">
      <ResponsePanel :response="response" :send-error="sendError" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { Button } from '@/shared/components/ui'
import { ParamSection, BodyEditor, ResponsePanel } from '@/features/debug-panel'
import { PlusIcon, XIcon } from '@/shared/components/icons'
import { getMethodColor } from '@/shared/constants/httpMethods'
import { makeSchemaHelpers } from '@/shared/utils/schema'
import { useLocale } from '@/shared/hooks'
import { useSettings } from '@/features/settings'
import { createHistory } from '@/shared/utils/history'
import { startDragResize } from '@/shared/hooks'

const props = defineProps({
  endpoint: { type: Object, required: true },
  spec: { type: Object, default: null },
  baseUrl: { type: String, default: '' },
  basePath: { type: String, default: '' },
  mobile: { type: Boolean, default: false },
  enableHistory: { type: Boolean, default: true },
})

defineEmits(['close', 'drag-start'])

const { t } = useLocale()
const { applySettings } = useSettings()

const width = ref(380)
const height = ref(220)

defineExpose({ panelWidth: width })

function startResponseResize(e) {
  startDragResize(e, height, { axis: 'y', min: 80, max: 600, sign: -1 })
}

function startResize(e) {
  startDragResize(e, width, { axis: 'x', min: 280, max: 600, sign: -1 })
}

const hasBody = computed(() => ['post', 'put', 'patch'].includes(props.endpoint?.method?.toLowerCase()))

const path = ref([])

const resolved = computed(() => {
  let p = props.endpoint?.path ?? ''
  for (const param of path.value) {
    if (param.value) p = p.replace(`{${param.name}}`, encodeURIComponent(param.value))
  }
  return p
})

const url = computed(() => props.baseUrl + props.basePath + resolved.value)

const query = ref([])
function addQueryRow() {
  query.value.push({ name: '', value: '', enabled: true, meta: null, custom: true })
}

const isQueryVisible = ref(true)
const isHeadersVisible = ref(false)
const isCookiesVisible = ref(false)
const isBodyVisible = ref(false)
const isHistoryVisible = ref(false)
const headers = ref([])
const cookies = ref([])
const body = ref('')
const type = ref('json')
const bodyEditorRef = ref()

const response = ref(null)
const sendError = ref('')
const sending = ref(false)

const CT = {
  json: 'application/json',
  xml: 'application/xml',
  text: 'text/plain',
  graphql: 'application/json',
}

async function sendRequest() {
  if (sending.value) return
  response.value = null
  sendError.value = ''
  sending.value = true
  try {
    let u = url.value
    const searchParams = new URLSearchParams()
    for (const r of query.value) {
      if (r.enabled && r.name) searchParams.set(r.name, r.value)
    }
    const method = props.endpoint.method.toUpperCase()
    const reqHeaders = {}
    const hasBodySection = hasBody.value || isBodyVisible.value
    const btype = type.value

    for (const h of headers.value) {
      if (h.enabled && h.name) reqHeaders[h.name] = h.value
    }
    const cookieStr = cookies.value
      .filter((c) => c.enabled && c.name)
      .map((c) => `${c.name}=${c.value}`)
      .join('; ')
    if (cookieStr) reqHeaders['Cookie'] = cookieStr
    applySettings(reqHeaders, searchParams)
    const ps = searchParams.toString()
    if (ps) u += '?' + ps

    const init = { method, headers: reqHeaders }
    if (hasBodySection && btype !== 'none') {
      if (CT[btype]) {
        reqHeaders['Content-Type'] = CT[btype]
        if (body.value.trim()) init.body = body.value
      } else if (btype === 'form-urlencoded') {
        reqHeaders['Content-Type'] = 'application/x-www-form-urlencoded'
        const params = new URLSearchParams()
        for (const r of bodyEditorRef.value?.formRows ?? []) {
          if (r.name) params.set(r.name, r.value)
        }
        init.body = params.toString()
      } else if (btype === 'form-data') {
        const fd = new FormData()
        for (const r of bodyEditorRef.value?.formRows ?? []) {
          if (!r.name) continue
          if (r.type === 'file' && r.file) fd.append(r.name, r.file, r.file.name)
          else if (r.type !== 'file') fd.append(r.name, r.value)
        }
        init.body = fd
      }
    }

    const t0 = Date.now()
    const res = await fetch(u, init)
    const ms = Date.now() - t0
    const rawText = await res.text()
    let responseBody = rawText
    try {
      responseBody = JSON.stringify(JSON.parse(rawText), null, 2)
    } catch {}
    const bytes = new TextEncoder().encode(rawText).length
    const size = bytes > 1024 ? `${(bytes / 1024).toFixed(1)} KB` : `${bytes} B`
    response.value = { status: res.status, time: ms, size, body: responseBody }
    createHistory(props.endpoint.id).push({
      request: {
        method: props.endpoint.method,
        url: url.value,
        headers: Object.fromEntries(headers.value.filter((h) => h.enabled).map((h) => [h.name, h.value])),
        params: Object.fromEntries(query.value.filter((q) => q.enabled).map((q) => [q.name, q.value])),
        body: body.value,
        bodyType: btype,
      },
      response: { status: res.status, statusText: res.statusText, headers: {}, data: responseBody, time: ms },
    })
    historyItems.value = createHistory(props.endpoint.id).list()
  } catch (e) {
    sendError.value = e.message
  } finally {
    sending.value = false
  }
}

const historyItems = ref([])

function loadHistory() {
  historyItems.value = createHistory(props.endpoint?.id ?? '').list()
}

function onOpenHistory() {
  isHistoryVisible.value = true
  loadHistory()
}

function restoreHistory(item) {
  query.value = item.queryRows ?? []
  headers.value = item.headerRows ?? []
  cookies.value = item.cookieRows ?? []
  path.value = item.pathParams ?? []
  if (bodyEditorRef.value) {
    bodyEditorRef.value.reset(item.bodyType ?? 'json', item.body ?? '', item.bodyFormRows ?? [])
  } else {
    body.value = item.body ?? ''
    type.value = item.bodyType ?? 'json'
  }
  isHistoryVisible.value = false
}

function onClearHistory() {
  createHistory(props.endpoint?.id ?? '').clear()
  historyItems.value = []
}

function fmtTime(iso) {
  try {
    const d = new Date(iso)
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  } catch {
    return iso
  }
}

watch(
  () => props.endpoint,
  (ep) => {
    response.value = null
    sendError.value = ''
    const isBodyMethod = ['post', 'put', 'patch'].includes(ep?.method?.toLowerCase())
    isQueryVisible.value = !isBodyMethod
    isHeadersVisible.value = false
    isCookiesVisible.value = false
    isBodyVisible.value = false
    isHistoryVisible.value = false
    historyItems.value = ep ? createHistory(ep.id).list() : []
    if (!ep) return
    path.value = (ep.parameters ?? [])
      .filter((p) => p.in === 'path')
      .map((p) => ({
        name: p.name,
        value: '',
        enabled: true,
        meta: {
          type: p.schema?.type ?? p.type ?? 'string',
          required: p.required ?? true,
          description: p.description ?? '',
        },
      }))
    query.value = (ep.parameters ?? [])
      .filter((p) => p.in === 'query')
      .map((p) => ({
        name: p.name,
        value: '',
        enabled: false,
        meta: {
          type: p.schema?.type ?? p.type ?? 'string',
          required: p.required ?? false,
          description: p.description ?? '',
        },
        custom: false,
      }))
    // Check for Swagger 2.0 formData parameters
    const formDataParams = (ep.parameters ?? []).filter((p) => p.in === 'formData')

    if (formDataParams.length > 0) {
      // Swagger 2.0 style formData parameters
      const consumes = ep.consumes ?? []
      const formType = consumes.includes('multipart/form-data') ? 'formdata' : 'form'

      // Set type first so the Body editor shows the correct type
      type.value = formType
      body.value = ''

      // Initialize form rows with default values from formData parameters
      const formRows = formDataParams.map((p) => ({
        name: p.name,
        value: p.default !== undefined ? String(p.default) : '',
        type: 'text',
        file: null,
      }))

      // Use nextTick to ensure bodyEditorRef is initialized
      nextTick(() => {
        if (bodyEditorRef.value) {
          bodyEditorRef.value.reset(formType, '', formRows)
        }
      })
    } else if (hasBody.value && ep.requestBody) {
      const content = ep.requestBody.content ?? {}
      if (content['application/x-www-form-urlencoded']) type.value = 'form-urlencoded'
      else if (content['multipart/form-data']) type.value = 'form-data'
      else if (content['application/xml'] || content['text/xml']) type.value = 'xml'
      else if (content['text/plain']) type.value = 'text'
      else type.value = 'json'

      const { resolveRef, schemaToExample, schemaToXml } = makeSchemaHelpers(props.spec)
      const media =
        content['application/json'] ?? content['application/xml'] ?? content['text/xml'] ?? Object.values(content)[0]

      console.log('Body type:', type.value)
      console.log('Content keys:', Object.keys(content))
      console.log('Media:', media)

      if (media?.schema) {
        try {
          const resolved = resolveRef(media.schema)
          console.log('Resolved schema:', resolved)

          // Generate appropriate default based on content type
          if (type.value === 'xml') {
            const xmlContent = schemaToXml(resolved, 'XMLRequest')
            console.log('Generated XML:', xmlContent)
            console.log('XML length:', xmlContent.length)
            body.value = xmlContent
          } else if (type.value === 'json') {
            const jsonContent = JSON.stringify(schemaToExample(resolved), null, 2)
            console.log('Generated JSON:', jsonContent)
            body.value = jsonContent
          } else if (type.value === 'text') {
            body.value = ''
          } else {
            body.value = ''
          }
          console.log('Final body.value:', body.value)
        } catch (err) {
          console.error('Error generating body:', err)
          // Fallback to empty string for non-JSON types
          body.value = type.value === 'json' ? '{}' : ''
        }
      } else {
        // No schema, use appropriate empty default
        body.value = type.value === 'json' ? '{}' : ''
      }
    } else {
      type.value = 'json'
      body.value = ''
    }
  },
  { immediate: true }
)
</script>
