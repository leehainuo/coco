<template>
  <component
    :is="as"
    v-bind="{ ...linkAttrs, ...$attrs }"
    :disabled="as === 'button' ? disabled || loading : undefined"
    :class="classes"
  >
    <svg v-if="loading" class="h-3.5 w-3.5 shrink-0 animate-spin" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
    </svg>
    <slot />
  </component>
</template>

<script setup>
/**
 * Button - Button component
 * Supports multiple styles, sizes, loading states, can be rendered as button or a tag
 * @component
 * @example
 * <Button variant="primary" size="md" :loading="true">Submit</Button>
 */
import { computed, useAttrs } from 'vue'
import { twMerge } from 'tailwind-merge'

const props = defineProps({
  /** @type {string} Style variant: primary, danger, outline, ghost */
  variant: { type: String, default: 'primary' },
  /** @type {string} Size: xs, sm, md */
  size: { type: String, default: 'md' },
  /** @type {boolean} Whether disabled */
  disabled: { type: Boolean, default: false },
  /** @type {boolean} Whether loading */
  loading: { type: Boolean, default: false },
  /** @type {string} HTML tag to render: button, a */
  as: { type: String, default: 'button' },
  /** @type {string} Link URL (when as="a") */
  href: { type: String, default: undefined },
  /** @type {string} Download filename (when as="a") */
  download: { type: String, default: undefined },
})

defineOptions({ inheritAttrs: false })

const attrs = useAttrs()

const linkAttrs = computed(() => (props.as === 'a' ? { href: props.href, download: props.download } : {}))

const variantMap = {
  primary: 'bg-base hover:bg-base/80 text-base-foreground',
  danger: 'bg-destructive hover:bg-destructive/80 text-base-foreground',
  outline: 'border border-border text-muted-foreground hover:bg-accent',
  ghost: 'text-muted-foreground hover:bg-accent',
}

const sizeMap = {
  xs: 'px-2 py-1 text-[11px] gap-1',
  sm: 'px-3 py-1.5 text-xs gap-1.5',
  md: 'px-4 py-2.5 text-sm gap-1.5',
}

// twMerge: external class wins over internal defaults (e.g. px-5 overrides px-4)
const classes = computed(() =>
  twMerge(
    'inline-flex items-center justify-center font-medium rounded-lg transition-colors shrink-0 select-none',
    variantMap[props.variant] ?? variantMap.primary,
    sizeMap[props.size] ?? sizeMap.md,
    props.disabled || props.loading ? 'opacity-50 cursor-not-allowed pointer-events-none' : '',
    attrs.class
  )
)
</script>
