<template>
  <Teleport to="body">
    <Transition name="drawer-fade">
      <div v-if="mobileOpen" class="bg-overlay fixed inset-0 z-40 lg:hidden" @click="$emit('close-mobile')" />
    </Transition>
  </Teleport>

  <aside
    :style="isMobile ? {} : { width: width + 'px' }"
    :class="[
      'border-sidebar-border bg-sidebar fixed inset-y-0 left-0 z-50 flex h-full w-80 shrink-0 flex-col overflow-hidden border-r',
      'lg:relative lg:z-auto lg:w-auto',
      isMobile && !isResizing ? 'transition-transform duration-300 ease-in-out' : '',
      mobileOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
    ]"
  >
    <!-- Top -->
    <div class="flex shrink-0 items-center gap-2.5 px-5 pt-5 pb-2.5">
      <!-- API icon -->
      <!-- <div class="w-7 h-7 rounded-lg bg-blue-500 flex items-center justify-center shrink-0">
        <svg class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
        </svg>
      </div> -->

      <span class="text-sidebar-primary-foreground flex-1 truncate text-sm font-semibold">{{ title }}</span>
      <LangToggle :locale="locale" @toggle="$emit('toggle-lang')" />
      <ThemeToggle :dark="dark" @toggle="$emit('toggle-dark')" />
      <CloseToggle @toggle="$emit('close-mobile')" />
    </div>

    <!-- Search -->
    <SidebarSearch v-model="search" />

    <!-- Tag groups -->
    <nav class="flex-1 overflow-y-auto px-2.5 pb-4">
      <div v-for="tag in filtered" :key="tag.name">
        <!-- Tag row -->
        <button
          class="text-sidebar-primary-foreground hover:bg-sidebar-accent flex w-full items-center justify-between rounded-lg px-2.5 py-2.5 text-sm font-medium transition-colors"
          @click="onToggleTag(tag.name)"
        >
          <span class="truncate">{{ tag.name }}</span>
          <ChevronRightIcon v-if="!openTags.has(tag.name)" :size="14" color="var(--icon-foreground)" class="shrink-0" />
          <ChevronDownIcon v-else :size="14" color="var(--icon-foreground)" class="shrink-0" />
        </button>

        <!-- Endpoints -->
        <ul
          v-show="openTags.has(tag.name) || search.trim().length > 0"
          class="border-sidebar-border ml-4 border-l pb-2 pl-2"
        >
          <li v-for="ep in tag.endpoints" :key="ep.id" class="py-0.5">
            <button
              :class="[
                'flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left transition-all duration-150',
                selected?.id === ep.id ? 'bg-base-background' : 'hover:bg-sidebar-accent',
              ]"
              @click="$emit('select', ep)"
            >
              <span
                :class="[
                  'flex-1 truncate text-sm',
                  selected?.id === ep.id ? 'text-base font-medium' : 'text-sidebar-primary-foreground',
                ]"
              >
                {{ ep.summary || ep.path }}
              </span>
              <MethodBadge :method="ep.method" class="shrink-0 rounded-[5px] text-[0.55rem]" />
            </button>
          </li>
        </ul>
      </div>

      <p v-if="filtered.length === 0" class="text-sidebar-muted-foreground px-4 py-8 text-center text-xs">
        {{ t('noResults') }}
      </p>
    </nav>

    <!-- Global Settings menu item -->
    <div class="border-sidebar-border shrink-0 border-t px-2.5 py-5">
      <SettingsMenuItem
        :label="t('globalSettings')"
        :active="props.settingsOpen"
        :has-settings="hasSettings"
        @click="$emit('open-settings')"
      />
    </div>

    <!-- Resize handle (desktop only) -->
    <div @mousedown.prevent="startResize" class="resize-handle resize-handle-y right-0 hidden lg:block" />
  </aside>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ThemeToggle, LangToggle, SidebarSearch, CloseToggle, SettingsMenuItem } from '@/features/sidebar'
import { MethodBadge } from '@/shared/components/ui'
import { ChevronRightIcon, ChevronDownIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'
import { useSettings } from '@/features/settings'
import { startDragResize } from '@/shared/hooks'

const props = defineProps({
  tags: { type: Array, default: () => [] },
  selected: { type: Object, default: null },
  title: { type: String, default: 'API Docs' },
  version: { type: String, default: '' },
  dark: { type: Boolean, default: false },
  locale: { type: String, default: 'zh' },
  mobileOpen: { type: Boolean, default: false },
  settingsOpen: { type: Boolean, default: false },
})

defineEmits(['select', 'toggle-dark', 'toggle-lang', 'close-mobile', 'open-settings'])

const { t } = useLocale()
const { settings } = useSettings()

const hasSettings = computed(() => settings.value.type !== 'none')

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

const search = ref('')
const openTags = ref(new Set())
const width = ref(312)

function startResize(e) {
  startDragResize(e, width, { axis: 'x', min: 160, max: 400, sign: 1 })
}

function onToggleTag(name) {
  const s = new Set(openTags.value)
  s.has(name) ? s.delete(name) : s.add(name)
  openTags.value = s
}

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return props.tags
  return props.tags
    .map((t) => {
      if (t.name.toLowerCase().includes(q)) return t
      const endpoints = t.endpoints.filter((ep) => (ep.summary || '').toLowerCase().includes(q))
      return endpoints.length ? { ...t, endpoints } : null
    })
    .filter(Boolean)
})
</script>
