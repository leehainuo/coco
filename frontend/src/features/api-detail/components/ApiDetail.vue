<template>
  <!-- Main wrapper -->
  <div v-if="endpoint" class="relative flex flex-1 overflow-hidden">
    <!-- Doc area: padding-right keeps content clear of the overlay debug panel -->
    <div
      class="flex min-w-0 flex-1 flex-col overflow-y-auto"
      :style="isDebugOpen && !isMobile ? { paddingRight: padding } : {}"
    >
      <div class="mx-auto flex w-full flex-1 flex-col px-10 py-8">
        <!-- Header: Breadcrumb + Title -->
        <ApiHeader :tag="(endpoint.tags && endpoint.tags[0]) || ''" :title="endpoint.summary || endpoint.path" />

        <!-- URL bar -->
        <UrlBar
          :method="endpoint.method"
          :base-url="specBaseUrl"
          :base-path="specBasePath"
          :path="endpoint.path"
          :debug-open="isDebugOpen"
          :enable-debug="props.config.enableDebug"
          @open-debug="isDebugOpen = true"
        />

        <!-- Request Parameters -->
        <ApiParams :params="params" :request-body="requestBody" :spec="spec" :consumes="consumes" />

        <!-- Response -->
        <ApiResponses :responses="responses" :spec="spec" :produces="produces" />

        <!-- Prev / Next -->
        <EndpointNav :prev="prevEndpoint" :next="nextEndpoint" @select="$emit('select', $event)" />

        <!-- Footer -->
        <DocFooter :swagger-json-path="openapiJsonPath" :enable-export="props.config.enableExport" />
      </div>
    </div>

    <!-- Debug panel: desktop overlay card -->
    <div v-if="isDebugOpen && !isMobile" class="pointer-events-none absolute inset-y-0 right-0 z-10 flex items-stretch">
      <DebugPanel
        ref="panel"
        class="pointer-events-auto"
        :endpoint="endpoint"
        :spec="spec"
        :base-url="specBaseUrl"
        :base-path="specBasePath"
        :enable-history="props.config.enableHistory"
        @close="isDebugOpen = false"
      />
    </div>

    <!-- Debug panel: mobile bottom sheet -->
    <Teleport to="body">
      <Transition name="drawer-fade" :css="!isResizing">
        <div v-if="isDebugOpen && isMobile" class="bg-overlay fixed inset-0 z-40" @click="isDebugOpen = false" />
      </Transition>
      <Transition name="sheet-slide" :css="!isResizing">
        <div
          v-if="isDebugOpen && isMobile"
          class="fixed inset-x-0 bottom-0 z-50 flex flex-col"
          :style="{ height: height + 'px' }"
        >
          <DebugPanel
            :mobile="true"
            :endpoint="endpoint"
            :spec="spec"
            :base-url="specBaseUrl"
            :base-path="specBasePath"
            :enable-history="props.config.enableHistory"
            @close="isDebugOpen = false"
            @drag-start="startSheetResize"
          />
        </div>
      </Transition>
    </Teleport>
  </div>

  <!-- Empty state -->
  <EmptyState v-else :icon="FileBracesIcon" :icon-size="56" :description="t('selectApi')" class="flex-1" />
</template>

<script setup>
/**
 * ApiDetail - API detail page component
 * Display complete API endpoint documentation including parameters, request body, responses, etc.
 * Supports debug panel for both desktop and mobile
 * @component
 */
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ApiHeader, ApiParams, ApiResponses, UrlBar, EndpointNav, DocFooter } from '@/features/api-detail'
import { DebugPanel } from '@/features/debug-panel'
import { EmptyState } from '@/shared/components'
import { FileBracesIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'

const props = defineProps({
  endpoint: { type: Object, default: null },
  spec: { type: Object, default: null },
  allEndpoints: { type: Array, default: () => [] },
  config: { type: Object, default: () => ({ enableDebug: true, enableExport: true, enableHistory: true }) },
})

defineEmits(['select'])

const { t } = useLocale()

const isDebugOpen = ref(false)
const panel = ref(null)

const isMobile = ref(window.innerWidth < 1024)
const isResizing = ref(false)
let resizeTimer = null
function onResize() {
  isResizing.value = true
  isMobile.value = window.innerWidth < 1024
  clearTimeout(resizeTimer)
  resizeTimer = setTimeout(() => {
    isResizing.value = false
  }, 200)
}
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => {
  window.removeEventListener('resize', onResize)
  clearTimeout(resizeTimer)
})

const height = ref(Math.round(window.innerHeight * 0.9))

function startSheetResize(e) {
  const startY = e.touches ? e.touches[0].clientY : e.clientY
  const startH = height.value
  function onMove(ev) {
    const y = ev.touches ? ev.touches[0].clientY : ev.clientY
    height.value = Math.min(Math.round(window.innerHeight * 0.95), Math.max(280, startH + (startY - y)))
  }
  function onUp() {
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
    document.removeEventListener('touchmove', onMove)
    document.removeEventListener('touchend', onUp)
  }
  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
  document.addEventListener('touchmove', onMove, { passive: false })
  document.addEventListener('touchend', onUp)
}

const padding = computed(() => {
  const w = panel.value?.panelWidth ?? 380
  return w + 24 + 'px'
})

const params = computed(() => props.endpoint?.parameters ?? [])
const requestBody = computed(() => props.endpoint?.requestBody ?? null)
const responses = computed(() => props.endpoint?.responses ?? null)
const consumes = computed(() => props.endpoint?.consumes ?? [])
const produces = computed(() => {
  const p = props.endpoint?.produces ?? []
  return p.length > 0 ? p[0] : 'application/json'
})

const currentIdx = computed(() => props.allEndpoints.findIndex((e) => e.id === props.endpoint?.id))
const prevEndpoint = computed(() => (currentIdx.value > 0 ? props.allEndpoints[currentIdx.value - 1] : null))
const nextEndpoint = computed(() =>
  currentIdx.value < props.allEndpoints.length - 1 ? props.allEndpoints[currentIdx.value + 1] : null
)

const specBaseUrl = computed(() => {
  if (props.spec?.servers?.[0]?.url) {
    try {
      return new URL(props.spec.servers[0].url).origin
    } catch {
      return ''
    }
  }
  return window.location.origin
})
const specBasePath = computed(() => {
  if (props.spec?.servers?.[0]?.url) {
    try {
      return new URL(props.spec.servers[0].url).pathname.replace(/\/$/, '')
    } catch {
      return ''
    }
  }
  return (props.spec?.basePath ?? '').replace(/\/$/, '')
})

const openapiJsonPath = computed(() => {
  const base = window.location.pathname.replace(/\/+$/, '')
  return base + '/openapi.json'
})
</script>
