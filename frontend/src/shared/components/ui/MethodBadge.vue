<template>
  <span :class="['inline-flex items-center px-2 py-0.5 font-bold tracking-wide uppercase', colorClass]">
    {{ method }}
  </span>
</template>

<script setup>
/**
 * MethodBadge - HTTP method badge component
 * Display HTTP methods (GET, POST, PUT, DELETE, etc.) with auto-applied colors
 * @component
 * @example
 * <MethodBadge method="GET" />
 */
import { computed } from 'vue'
import { getMethodColor } from '@/shared/constants/httpMethods'

const props = defineProps({
  /** @type {string} HTTP method name */
  method: {
    type: String,
    required: true,
  },
})

const colorClass = computed(() => getMethodColor(props.method))
</script>
