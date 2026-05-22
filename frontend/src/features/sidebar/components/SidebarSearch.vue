<template>
  <div class="shrink-0 px-5 pt-3 pb-2">
    <div class="border-sidebar-border bg-sidebar-muted flex items-center gap-2 rounded-lg border px-2.5 py-2">
      <SearchIcon size="16" color="var(--sidebar-primary-foreground)" class="shrink-0" />
      <input
        :value="modelValue"
        type="text"
        :placeholder="t('search')"
        class="text-sidebar-primary-foreground placeholder-sidebar-muted-foreground min-w-0 flex-1 bg-transparent text-xs focus:outline-none"
        @input="onInput"
      />
    </div>
  </div>
</template>

<script setup>
import { useLocale } from '@/shared/hooks'
import { SearchIcon } from '@/shared/components/icons'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue'])

const { t } = useLocale()

function onInput(event) {
  emit('update:modelValue', event.target.value)
}
</script>
