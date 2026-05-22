<template>
  <div class="border-border flex shrink-0 items-center justify-between border-b px-4 py-3">
    <span class="text-primary-foreground text-sm font-semibold">{{ t('onlineDebug') }}</span>
    <button
      @click="$emit('close')"
      class="text-icon-foreground hover:bg-icon-accent hover:text-icon-accent-foreground rounded p-1 transition-colors"
    >
      <XIcon :size="16" />
    </button>
  </div>
</template>

<script setup>
import { XIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'

const { t } = useLocale()

defineEmits(['close'])
</script>
