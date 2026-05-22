<template>
  <div class="flex-1 overflow-y-auto">
    <div class="mx-auto flex w-full max-w-2xl flex-col px-10 py-8">
      <!-- Title -->
      <div class="mb-1 flex items-center justify-between">
        <h1 class="text-primary-foreground text-2xl font-bold">{{ t('settings') }}</h1>
        <button
          @click="clearSettings"
          :title="t('settingsClear')"
          class="text-icon-foreground hover:bg-icon-accent hover:text-destructive rounded-lg p-1 transition-colors"
        >
          <BrushCleaningIcon :size="16" />
        </button>
      </div>
      <p class="text-muted-foreground mb-8 text-sm">{{ t('settingsPageDesc') }}</p>

      <!-- Type selector -->
      <section class="mb-8">
        <h2 class="text-primary-foreground mb-3 font-semibold">{{ t('settingsType') }}</h2>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="opt in types"
            :key="opt.value"
            :class="[
              'rounded-lg border px-3.5 py-1.5 text-sm font-medium transition-colors',
              settings.type === opt.value
                ? 'border-base text-base'
                : 'border-border text-primary-foreground hover:border-base hover:text-base',
            ]"
            @click="settings.type = opt.value"
          >
            {{ opt.label }}
          </button>
        </div>
      </section>

      <!-- None -->
      <section v-if="settings.type === 'none'" class="border-border rounded-xl border px-6 py-5">
        <p class="text-muted-foreground text-sm">{{ t('settingsNoneHint') }}</p>
      </section>

      <!-- Bearer Token -->
      <section v-if="settings.type === 'bearer'" class="border-border space-y-4 rounded-xl border px-6 py-5">
        <div>
          <label class="text-primary-foreground mb-1.5 block text-sm font-medium">Token</label>
          <input
            v-model="settings.bearer.token"
            type="password"
            placeholder="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
            class="border-border bg-muted text-primary-foreground focus:border-base w-full rounded-lg border px-3 py-2 font-mono text-sm transition-colors focus:outline-none"
          />
        </div>
        <p class="text-muted-foreground text-xs">{{ t('settingsBearerHint') }}</p>
      </section>

      <!-- API Key -->
      <section v-if="settings.type === 'apikey'" class="border-border space-y-4 rounded-xl border px-6 py-5">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="text-primary-foreground mb-1.5 block text-sm font-medium">{{ t('settingsKeyName') }}</label>
            <input
              v-model="settings.apikey.name"
              type="text"
              placeholder="X-API-Key"
              class="border-border bg-muted text-primary-foreground focus:border-base w-full rounded-lg border px-3 py-2 font-mono text-sm transition-colors focus:outline-none"
            />
          </div>
          <div>
            <label class="text-primary-foreground mb-1.5 block text-sm font-medium">{{ t('settingsKeyValue') }}</label>
            <input
              v-model="settings.apikey.value"
              type="password"
              class="border-border bg-muted text-primary-foreground focus:border-base w-full rounded-lg border px-3 py-2 font-mono text-sm transition-colors focus:outline-none"
            />
          </div>
        </div>
        <div>
          <label class="text-primary-foreground mb-2 block text-sm font-medium">{{ t('settingsKeyIn') }}</label>
          <div class="flex gap-5">
            <label class="flex cursor-pointer items-center gap-2">
              <input type="radio" v-model="settings.apikey.in" value="header" class="accent-base" />
              <span class="text-muted-foreground text-sm">Header</span>
            </label>
            <label class="flex cursor-pointer items-center gap-2">
              <input type="radio" v-model="settings.apikey.in" value="query" class="accent-base" />
              <span class="text-muted-foreground text-sm">Query Param</span>
            </label>
          </div>
        </div>
      </section>

      <!-- Basic Auth -->
      <section v-if="settings.type === 'basic'" class="border-border space-y-4 rounded-xl border px-6 py-5">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="text-primary-foreground mb-1.5 block text-sm font-medium">{{ t('settingsUsername') }}</label>
            <input
              v-model="settings.basic.username"
              type="text"
              class="text-primary-foreground border-border bg-muted focus:border-base w-full rounded-lg border px-3 py-2 font-mono text-sm transition-colors focus:outline-none"
            />
          </div>
          <div>
            <label class="text-primary-foreground mb-1.5 block text-sm font-medium">{{ t('settingsPassword') }}</label>
            <input
              v-model="settings.basic.password"
              type="password"
              class="text-primary-foreground border-border bg-muted focus:border-base w-full rounded-lg border px-3 py-2 font-mono text-sm transition-colors focus:outline-none"
            />
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useSettings } from '@/features/settings'
import { useLocale } from '@/shared/hooks'
import { BrushCleaningIcon } from '@/shared/components/icons'

const { settings } = useSettings()
const { t } = useLocale()

const types = computed(() => [
  { value: 'none', label: t('settingsNone') },
  { value: 'bearer', label: 'Bearer Token' },
  { value: 'apikey', label: 'API Key' },
  { value: 'basic', label: 'Basic Auth' },
])

function clearSettings() {
  settings.value.type = 'none'
  settings.value.bearer.token = ''
  settings.value.apikey.value = ''
  settings.value.basic.username = ''
  settings.value.basic.password = ''
}
</script>
