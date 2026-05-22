<template>
  <button
    @click="$emit('toggle')"
    class="text-icon-foreground hover:text-icon-accent-foreground hover:bg-icon-accent shrink-0 rounded-md p-1 transition-colors"
  >
    <SunIcon v-if="dark" :size="16" color="currentColor" />
    <MoonIcon v-else :size="16" color="currentColor" />
  </button>
</template>

<script setup>
import { SunIcon, MoonIcon } from '@/shared/components/icons'

defineProps({
  dark: { type: Boolean, default: false },
})

defineEmits(['toggle'])
</script>
