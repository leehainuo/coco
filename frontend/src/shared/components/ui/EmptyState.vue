<template>
  <div :class="['flex flex-col items-center justify-center gap-3 py-12', className]">
    <div v-if="icon" class="text-muted-foreground opacity-40">
      <component :is="icon" :size="iconSize" />
    </div>
    <div class="text-center">
      <p v-if="title" class="text-primary-foreground mb-1 font-medium">{{ title }}</p>
      <p v-if="description" class="text-muted-foreground text-sm">{{ description }}</p>
    </div>
    <div v-if="$slots.action">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup>
/**
 * EmptyState - Empty state placeholder component
 * Display empty states like no data, no results, etc.
 * @component
 * @example
 * <EmptyState
 *   :icon="SearchIcon"
 *   title="No results found"
 *   description="Try different search keywords"
 * >
 *   <template #action>
 *     <button>Reset Search</button>
 *   </template>
 * </EmptyState>
 */
defineProps({
  /** @type {Object} Icon component */
  icon: {
    type: Object,
    default: null,
  },
  /** @type {number} Icon size (pixels) */
  iconSize: {
    type: Number,
    default: 48,
  },
  /** @type {string} Title text */
  title: {
    type: String,
    default: '',
  },
  /** @type {string} Description text */
  description: {
    type: String,
    default: '',
  },
  /** @type {string} Additional CSS class */
  className: {
    type: String,
    default: '',
  },
})
</script>
