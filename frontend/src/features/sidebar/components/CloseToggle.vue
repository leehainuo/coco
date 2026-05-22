<template>
  <button
    @click="$emit('toggle')"
    class="text-icon-foreground hover:text-icon-accent-foreground hover:bg-icon-accent -mr-1 shrink-0 rounded-md p-1 transition-colors lg:hidden"
  >
    <XIcon :size="16" color="currentColor" />
  </button>
</template>

<script setup>
import { XIcon } from '@/shared/components/icons'

defineEmits(['toggle'])
</script>
