<template>
  <div class="bg-background border-border mb-4 flex items-stretch rounded-xl border px-2 py-2">
    <div class="flex shrink-0 items-center">
      <span
        :class="['inline-flex h-6 items-center rounded-md p-2 text-xs font-semibold uppercase', getMethodColor(method)]"
      >
        {{ method }}
      </span>
    </div>

    <div class="flex min-w-0 flex-1 cursor-pointer items-center px-2 select-none" @click="copy">
      <div class="group relative w-fit min-w-0">
        <span
          class="text-primary-foreground font-mono text-xs break-all group-hover:underline group-hover:decoration-dashed"
        >
          {{ baseUrl }}{{ basePath }}{{ path }}
        </span>
        <span
          class="text-primary-foreground bg-background border-border pointer-events-none absolute -top-6 left-1/2 z-10 -translate-x-1/2 rounded-md border px-2 py-1 text-xs whitespace-nowrap opacity-0 transition-opacity duration-150 group-hover:opacity-100"
        >
          {{ copied ? t('copied') : t('clickToCopy') }}
        </span>
      </div>
    </div>

    <Button
      v-if="enableDebug && !debugOpen"
      variant="debug"
      class="self-center rounded-md px-3 py-1"
      @click="$emit('openDebug')"
    >
      <PlayIcon :size="12" class="shrink-0" />
      <span>{{ t('debug') }}</span>
    </Button>
  </div>
</template>

<script setup>
import { Button } from '@/shared/components/ui'
import { PlayIcon } from '@/shared/components/icons'
import { getMethodColor } from '@/shared/constants/httpMethods'
import { useLocale, useClipboard } from '@/shared/hooks'

const props = defineProps({
  method: { type: String, required: true },
  baseUrl: { type: String, default: '' },
  basePath: { type: String, default: '' },
  path: { type: String, default: '' },
  debugOpen: { type: Boolean, default: false },
  enableDebug: { type: Boolean, default: true },
})

defineEmits(['openDebug'])

const { t } = useLocale()
const { copied, copy: copyToClipboard } = useClipboard()

function copy() {
  copyToClipboard(props.baseUrl + props.basePath + props.path)
}
</script>
