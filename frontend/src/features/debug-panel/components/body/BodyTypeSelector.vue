<template>
  <Dropdown placement="bottom-end" content-class="w-fit max-w-56">
    <template #trigger="{ open }">
      <button
        class="text-muted-foreground hover:bg-accent flex items-center gap-1 rounded px-2 py-0.5 text-xs transition-colors"
      >
        {{ currentLabel }}
        <ChevronDownIcon :size="12" class="opacity-60" />
      </button>
    </template>

    <template #default="{ close }">
      <div class="flex flex-col gap-1 p-2">
        <button
          v-for="opt in BODY_TYPES"
          :key="opt.value"
          @click="selectType(opt.value, close)"
          :class="[
            'w-full rounded-md px-4 py-2 text-left text-xs transition-colors',
            opt.value === modelValue
              ? 'text-primary-foreground bg-muted font-medium'
              : 'text-primary-foreground hover:bg-accent',
          ]"
        >
          {{ opt.label }}
        </button>
      </div>
    </template>
  </Dropdown>
</template>

<script setup>
/**
 * BodyTypeSelector - Request body type selector
 * Use dropdown menu to select request body type (None, JSON, XML, Form, etc.)
 * @component
 * @example
 * <BodyTypeSelector v-model="bodyType" />
 */
import { computed } from 'vue'
import { Dropdown } from '@/shared/components/ui'
import { ChevronDownIcon } from '@/shared/components/icons'
import { BODY_TYPES } from '@/shared/constants/bodyTypes'

const props = defineProps({
  /** @type {string} Currently selected request body type */
  modelValue: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['update:modelValue'])

const currentLabel = computed(() => {
  const type = BODY_TYPES.find((t) => t.value === props.modelValue)
  return type?.label || 'None'
})

function selectType(value, close) {
  emit('update:modelValue', value)
  close()
}
</script>
