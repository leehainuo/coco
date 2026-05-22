<template>
  <button
    :class="[
      'flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors',
      active ? 'bg-base-background text-base' : 'text-sidebar-primary-foreground hover:bg-sidebar-accent',
    ]"
    @click="$emit('click')"
  >
    <FileSlidersIcon :size="16" class="shrink-0" />
    <span class="flex-1 text-left font-medium">{{ label }}</span>
    <span v-if="hasSettings" class="bg-base h-1.5 w-1.5 shrink-0 rounded-full" />
  </button>
</template>

<script setup>
import { FileSlidersIcon } from '@/shared/components/icons'

defineProps({
  label: { type: String, required: true },
  active: { type: Boolean, default: false },
  hasSettings: { type: Boolean, default: false },
})

defineEmits(['click'])
</script>
