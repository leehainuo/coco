<template>
  <div class="flex h-full flex-col overflow-hidden">
    <div class="flex shrink-0 items-center justify-between px-3 py-2">
      <div class="flex items-center gap-1.5">
        <span class="text-primary-foreground text-xs font-medium">{{ t('result') }}</span>
        <template v-if="response">
          <span class="text-muted-foreground text-[10px]">{{ response.time }}ms · {{ response.size }}</span>
        </template>
      </div>
      <StatusDot v-if="response" :code="response.status" />
    </div>

    <div v-if="sendError" class="border-border bg-red-background mx-3 mb-2 shrink-0 rounded border">
      <p class="text-destructive text-xs">{{ sendError }}</p>
    </div>

    <div
      v-if="!response && !sendError"
      class="text-muted-foreground flex flex-1 flex-col items-center justify-center gap-3 select-none"
    >
      <RocketIcon :size="40" class="opacity-40" />
      <p class="text-xs">{{ t('sendHint') }}</p>
    </div>

    <div v-if="response" class="flex-1 overflow-y-auto px-3 py-2">
      <CodeBlock :code="response.body ?? ''" />
    </div>
  </div>
</template>

<script setup>
import { StatusDot, CodeBlock } from '@/shared/components/ui'
import { RocketIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'

defineProps({
  response: { type: Object, default: null },
  sendError: { type: String, default: '' },
})

const { t } = useLocale()
</script>
