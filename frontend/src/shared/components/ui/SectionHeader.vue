<template>
  <div class="mb-4 flex items-center justify-between">
    <h2 class="text-primary-foreground font-semibold">
      <slot />
    </h2>
    <div v-if="$slots.actions" class="flex items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
/**
 * SectionHeader - Section header component
 * Provides layout container for title and action buttons
 * @component
 * @example
 * <SectionHeader>
 *   Title
 *   <template #actions>
 *     <button>Action</button>
 *   </template>
 * </SectionHeader>
 */
</script>
