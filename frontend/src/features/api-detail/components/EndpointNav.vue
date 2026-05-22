<template>
  <div class="border-border flex items-center justify-between border-t py-5">
    <button
      v-if="prev"
      class="group text-icon-foreground hover:text-accent-foreground flex items-center gap-2 text-sm"
      @click="$emit('select', prev)"
    >
      <ChevronLeftIcon :size="16" class="transition-transform group-hover:-translate-x-1" />
      <span class="text-primary-foreground">
        <span class="text-muted-foreground block text-[10px]">{{ t('prevPage') }}</span>
        {{ prev.summary || prev.path }}
      </span>
    </button>
    <div v-else />

    <button
      v-if="next"
      class="group text-icon-foreground hover:text-accent-foreground flex items-center gap-2 text-right text-sm"
      @click="$emit('select', next)"
    >
      <span class="text-primary-foreground">
        <span class="text-muted-foreground block text-[10px]">{{ t('nextPage') }}</span>
        {{ next.summary || next.path }}
      </span>
      <ChevronRightIcon :size="16" class="transition-transform group-hover:translate-x-1" />
    </button>
    <div v-else />
  </div>
</template>

<script setup>
import { ChevronLeftIcon, ChevronRightIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'

defineProps({
  prev: { type: Object, default: null },
  next: { type: Object, default: null },
})

defineEmits(['select'])

const { t } = useLocale()
</script>
