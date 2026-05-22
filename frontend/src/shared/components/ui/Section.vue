<template>
  <section :class="['mb-8', className]">
    <slot />
  </section>
</template>

<script setup>
/**
 * Section - Section container component
 * Provides unified section spacing and layout
 * @component
 * @example
 * <Section className="custom-class">Content</Section>
 */
defineProps({
  /** @type {string} Additional CSS class */
  className: {
    type: String,
    default: '',
  },
})
</script>
