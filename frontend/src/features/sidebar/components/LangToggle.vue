<template>
  <Dropdown placement="bottom-end" content-class="min-w-32 py-1">
    <template #trigger>
      <button
        class="text-icon-foreground hover:text-icon-accent-foreground hover:bg-icon-accent shrink-0 rounded-md p-1 transition-colors"
        :title="t('language')"
      >
        <LanguagesIcon :size="16" color="currentColor" class="shrink-0" />
      </button>
    </template>

    <template #default="{ close }">
      <button
        v-for="option in options"
        :key="option.value"
        class="text-primary-foreground hover:bg-icon-accent flex w-full items-center justify-between gap-3 px-3 py-1.5 text-left text-sm transition-colors"
        @click="onSelect(option.value, close)"
      >
        <span class="truncate">{{ option.label }}</span>
        <CheckIcon v-if="locale === option.value" :size="14" color="currentColor" class="shrink-0" />
      </button>
    </template>
  </Dropdown>
</template>

<script setup>
import { computed } from 'vue'
import { Dropdown } from '@/shared/components/ui'
import { LanguagesIcon, CheckIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'

const props = defineProps({
  locale: { type: String, default: 'zh' },
  hasCustom: { type: Boolean, default: false },
  customName: { type: String, default: '' },
})

const emit = defineEmits(['select'])

const { t } = useLocale()

const options = computed(() => {
  const list = [
    { value: 'zh', label: '中文' },
    { value: 'en', label: 'English' },
  ]
  if (props.hasCustom) {
    list.push({ value: 'custom', label: props.customName || 'Custom' })
  }
  return list
})

function onSelect(value, close) {
  if (value !== props.locale) emit('select', value)
  close()
}
</script>
