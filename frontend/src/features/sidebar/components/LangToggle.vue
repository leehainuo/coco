<template>
  <button
    @click="$emit('toggle')"
    class="text-icon-foreground hover:text-icon-accent-foreground hover:bg-icon-accent shrink-0 rounded-md p-1 transition-colors"
    :title="locale === 'zh' ? 'Switch to English' : 'Switch to Chinese'"
  >
    <LanguagesIcon :size="16" color="currentColor" class="shrink-0" />
  </button>
</template>

<script setup>
import { LanguagesIcon } from '@/shared/components/icons'

defineProps({
  locale: { type: String, default: 'zh' },
})

defineEmits(['toggle'])
</script>
