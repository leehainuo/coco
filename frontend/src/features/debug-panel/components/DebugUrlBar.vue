<template>
  <div class="border-border flex shrink-0 items-start gap-2 border-b px-3 py-2.5">
    <div class="flex min-w-0 flex-1 items-start gap-1.5">
      <span :class="['shrink-0 rounded-md px-1.5 py-1 text-xs font-bold uppercase', getMethodColor(method)]">
        {{ method }}
      </span>
      <p class="text-primary-foreground min-w-0 flex-1 truncate pt-1 font-mono text-xs">
        {{ url }}
      </p>
    </div>
    <Button variant="primary" size="sm" :loading="loading" @click="$emit('send')">
      {{ loading ? t('sending') : t('send') }}
    </Button>
  </div>
</template>

<script setup>
import { Button } from '@/shared/components/ui'
import { useLocale } from '@/shared/hooks'
import { getMethodColor } from '@/shared/constants/httpMethods'

const { t } = useLocale()

const props = defineProps({
  method: {
    type: String,
    required: true,
  },
  url: {
    type: String,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['send'])
</script>
