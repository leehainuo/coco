<template>
  <div class="flex h-screen overflow-hidden font-sans antialiased">
    <!-- Mobile top bar -->
    <div class="fixed top-0 right-0 left-0 z-30 flex h-12 items-center px-3 lg:hidden">
      <button
        @click="isDrawerOpen = true"
        class="text-icon-foreground hover:bg-icon-accent hover:text-icon-accent-foreground shrink-0 rounded-md p-1.5 transition-colors"
      >
        <MenuIcon :size="18" color="currentColor" />
      </button>
      <span class="text-primary-foreground ml-2.5 truncate text-sm font-semibold">
        {{ spec?.info?.title ?? 'Coco API Docs' }}
      </span>
    </div>

    <!-- Sidebar -->
    <Sidebar
      :tags="tags"
      :selected="sel"
      :title="spec?.info?.title ?? 'Coco API Docs'"
      :version="spec?.info?.version ?? ''"
      :dark="dark"
      :locale="locale"
      :mobile-open="isDrawerOpen"
      :settings-open="isSettingsOpen"
      @select="onSelect"
      @toggle-dark="onToggleDark()"
      @close-mobile="isDrawerOpen = false"
      @toggle-lang="setLocale(locale === 'zh' ? 'en' : 'zh')"
      @open-settings="onOpenSettings"
    />

    <!-- Main content -->
    <main class="bg-background relative flex flex-1 overflow-hidden pt-12 lg:pt-0">
      <!-- Top gradient background -->
      <div
        class="pointer-events-none absolute inset-x-0 top-0 h-[460px] max-h-full w-full overflow-hidden"
        :style="{
          backgroundImage: `url(${dark ? darkBg : lightBg})`,
          backgroundSize: 'cover',
          backgroundPosition: 'top center',
          opacity: 1,
        }"
      />
      <SettingsPage v-if="isSettingsOpen" />
      <ApiDetail v-else :endpoint="sel" :spec="spec" :all-endpoints="all" :config="config" @select="onSelect" />
    </main>

    <!-- Loading -->
    <div v-if="loading" class="fixed top-0 right-0 left-0 z-50 h-1">
      <div class="animate-loading-bar h-full" style="background-color: var(--color-base)"></div>
    </div>

    <!-- Error -->
    <div v-if="err" class="bg-background absolute inset-0 z-50 flex items-center justify-center">
      <div class="max-w-sm px-6 text-center">
        <p class="text-destructive mb-2 font-medium">{{ t('loadError') }}</p>
        <p class="text-muted-foreground text-sm">{{ err }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watchEffect } from 'vue'
import { Sidebar } from '@/features/sidebar'
import { ApiDetail } from '@/features/api-detail'
import { SettingsPage } from '@/features/settings'
import { MenuIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'
import type { OpenAPISpec, Endpoint, Tag } from '@/shared/types'
import darkBg from '@/assets/dark-bg.webp'
import lightBg from '@/assets/light-bg.webp'

const { t, locale, setLocale } = useLocale()

const spec = ref<OpenAPISpec | null>(null)
const loading = ref(true)
const err = ref('')
const sel = ref<Endpoint | null>(null)
const isSettingsOpen = ref(false)
const isDrawerOpen = ref(false)
const config = ref({
  enableDebug: true,
  enableExport: true,
  enableHistory: true,
})
const saved = localStorage.getItem('coco:theme')
const dark = ref(saved === 'dark')

function onToggleDark() {
  dark.value = !dark.value
  localStorage.setItem('coco:theme', dark.value ? 'dark' : 'light')
}

watchEffect(() => {
  document.documentElement.classList.toggle('dark', dark.value)
})

function url(f: string): string {
  return window.location.pathname.replace(/\/+$/, '') + '/' + f
}

const METHODS = ['get', 'post', 'put', 'patch', 'delete', 'head', 'options']

const tags = computed<Tag[]>(() => {
  if (!spec.value?.paths) return []
  const m = new Map()
  for (const [p, item] of Object.entries(spec.value.paths)) {
    const pp: any = Array.isArray(item.parameters) ? item.parameters.filter((x: any) => !x.$ref) : []
    for (const [method, op] of Object.entries(item)) {
      if (!METHODS.includes(method)) continue
      const ops: any = (op.parameters ?? []).filter((x: any) => !x.$ref)
      const keys = new Set(ops.map((x: any) => `${x.in}:${x.name}`))
      const allParams = [...ops, ...pp.filter((x: any) => !keys.has(`${x.in}:${x.name}`))]

      // Separate body parameters from other parameters (Swagger 2.0 compatibility)
      const bodyParam = allParams.find((x: any) => x.in === 'body')
      const params = allParams.filter((x: any) => x.in !== 'body')

      // Convert Swagger 2.0 body parameter to requestBody format
      let requestBody = op.requestBody
      if (!requestBody && bodyParam) {
        // Get content type from consumes array
        const consumes = op.consumes || []
        const contentType = consumes[0] || 'application/json'

        requestBody = {
          required: bodyParam.required,
          description: bodyParam.description,
          content: {
            [contentType]: {
              schema: bodyParam.schema,
            },
          },
        }
      }

      const ts = op.tags?.length ? op.tags : ['Default']
      const ep = {
        id: `${method}:${p}`,
        method,
        path: p,
        summary: op.summary,
        description: op.description,
        parameters: params,
        requestBody: requestBody,
        responses: op.responses,
        consumes: op.consumes || [],
        produces: op.produces || [],
        tags: ts,
      }
      for (const t of ts) {
        if (!m.has(t)) m.set(t, [])
        m.get(t).push(ep)
      }
    }
  }
  return [...m.entries()].map(([name, endpoints]) => ({ name, endpoints }))
})

const all = computed<Endpoint[]>(() => {
  const seen = new Set()
  const list = []
  for (const t of tags.value) {
    for (const ep of t.endpoints) {
      if (!seen.has(ep.id)) {
        seen.add(ep.id)
        list.push(ep)
      }
    }
  }
  return list
})

onMounted(async () => {
  try {
    const [cfg, s] = await Promise.all([fetch(url('config.json')).catch(() => null), fetch(url('openapi.json'))])
    if (cfg?.ok) {
      try {
        const c = await cfg.json()
        if (c.title) document.title = c.title
        if (c.lang && !localStorage.getItem('coco:lang')) setLocale(c.lang, false)
        config.value.enableDebug = c.enableDebug ?? true
        config.value.enableExport = c.enableExport ?? true
        config.value.enableHistory = c.enableHistory ?? true
        if (!saved) {
          if (c.theme === 'dark') dark.value = true
          else if (c.theme === 'light') dark.value = false
          else dark.value = window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
        }
      } catch {
        if (!saved) dark.value = window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
      }
    } else if (!saved) {
      dark.value = window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
    }
    if (!s.ok) throw new Error(`HTTP ${s.status}`)
    spec.value = await s.json()
    if (!document.title || document.title === 'Coco API Docs') {
      document.title = spec.value?.info?.title ?? 'Coco API Docs'
    }
    const id = localStorage.getItem('coco:selected')
    if (id) {
      const ep = all.value.find((e) => e.id === id)
      if (ep) sel.value = ep
    }
  } catch (e) {
    err.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
})

function onSelect(ep: Endpoint): void {
  sel.value = ep
  isSettingsOpen.value = false
  isDrawerOpen.value = false
  localStorage.setItem('coco:selected', ep.id)
}

function onOpenSettings(): void {
  isSettingsOpen.value = true
  sel.value = null
  isDrawerOpen.value = false
  localStorage.removeItem('coco:selected')
}
</script>
