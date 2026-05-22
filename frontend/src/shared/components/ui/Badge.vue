<template>
  <span :class="classes">
    <span v-if="dot" :class="['h-1.5 w-1.5 shrink-0 rounded-full', dotColor]"></span>
    <slot />
  </span>
</template>

<script setup>
/**
 * Badge - Badge component
 * Display labels, status and other information with multiple colors and sizes
 * @component
 * @example
 * <Badge variant="blue" size="sm" :dot="true">Label</Badge>
 */
import { computed } from 'vue'

const props = defineProps({
  /** @type {string} Color variant: default, blue, green, red, orange, pink, gray */
  variant: { type: String, default: 'default' },
  /** @type {string} Size: xs, sm, md */
  size: { type: String, default: 'sm' },
  /** @type {boolean} Whether to show dot */
  dot: { type: Boolean, default: false },
})

const variantMap = {
  default: 'bg-muted text-muted-foreground',
  blue: 'bg-blue-background text-blue-foreground',
  green: 'bg-green-background text-green-foreground',
  red: 'bg-red-background text-red-foreground',
  orange: 'bg-orange-background text-orange-foreground',
  pink: 'bg-pink-background text-pink-foreground',
  gray: 'bg-muted text-muted-foreground',
}

const dotColorMap = {
  blue: 'bg-blue-foreground',
  green: 'bg-green-foreground',
  red: 'bg-red-foreground',
  orange: 'bg-orange-foreground',
  pink: 'bg-pink-foreground',
  gray: 'bg-muted-foreground',
  default: 'bg-muted-foreground',
}

const sizeMap = {
  xs: 'px-1.5 py-0.5 text-[10px] gap-1',
  sm: 'px-2 py-0.5 text-xs gap-1.5',
  md: 'px-2.5 py-1 text-sm gap-1.5',
}

const classes = computed(() => [
  'inline-flex items-center rounded-full font-medium',
  variantMap[props.variant] ?? variantMap.default,
  sizeMap[props.size] ?? sizeMap.sm,
])

const dotColor = computed(() => dotColorMap[props.variant] ?? dotColorMap.default)
</script>
