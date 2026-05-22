<template>
  <span class="inline-flex items-center gap-2">
    <span :class="['h-2 w-2 shrink-0 rounded-full', dotClass]"></span>
    <span :class="['font-mono text-sm font-semibold', textClass]">{{ code }}</span>
    <span v-if="description" class="text-primary-foreground text-sm">{{ description }}</span>
  </span>
</template>

<script setup>
/**
 * StatusDot - Status dot component
 * Display HTTP status code with corresponding color dot and description
 * @component
 * @example
 * <StatusDot :code="200" description="Success" />
 */
import { computed } from 'vue'

const props = defineProps({
  /** @type {string|number} HTTP status code */
  code: { type: [String, Number], required: true },
  /** @type {string} Status description */
  description: { type: String, default: '' },
})

function category(code) {
  const n = parseInt(code)
  if (n >= 200 && n < 300) return 'success'
  if (n >= 400 && n < 500) return 'warn'
  if (n >= 500) return 'error'
  return 'default'
}

const dotClass = computed(
  () =>
    ({
      success: 'bg-green-foreground',
      warn: 'bg-orange-foreground',
      error: 'bg-red-foreground',
      default: 'bg-muted-foreground',
    })[category(props.code)]
)

const textClass = computed(
  () =>
    ({
      success: 'text-green-foreground',
      warn: 'text-orange-foreground',
      error: 'text-red-foreground',
      default: 'text-muted-foreground',
    })[category(props.code)]
)
</script>
