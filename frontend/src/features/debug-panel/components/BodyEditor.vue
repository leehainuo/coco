<template>
  <section class="flex shrink-0 flex-col">
    <!-- Header: title + type selector -->
    <div class="border-border flex shrink-0 items-center justify-between px-3 py-2">
      <span class="text-primary-foreground text-xs font-medium">{{ t('body') }}</span>
      <BodyTypeSelector v-model="localBodyType" />
    </div>

    <!-- None state -->
    <EmptyState v-if="bodyType === 'none'" :description="t('bodyNone')" class="py-5" />

    <!-- Text editor: json / xml / text / binary / graphql -->
    <TextBodyEditor
      v-else-if="isTextType(bodyType)"
      :model-value="modelValue"
      :body-type="bodyType"
      :type-label="typeLabel"
      @update:model-value="$emit('update:modelValue', $event)"
    />

    <!-- Form editor: x-www-form-urlencoded / form-data -->
    <FormBodyEditor v-else-if="isFormType(bodyType)" v-model="formRows" :is-form-data="bodyType === 'formdata'" />
  </section>
</template>

<script setup>
/**
 * BodyEditor - Request body editor component
 * Supports multiple editor types: None, JSON, XML, Form, FormData, etc.
 * @component
 */
import { ref, computed, watch } from 'vue'
import { BodyTypeSelector, TextBodyEditor, FormBodyEditor } from '@/features/debug-panel'
import { EmptyState } from '@/shared/components/ui'
import { useLocale } from '@/shared/hooks'
import { BODY_TYPES, isTextType, isFormType } from '@/shared/constants'

const props = defineProps({
  modelValue: { type: String, default: '' },
  bodyType: { type: String, default: 'json' },
})

const emit = defineEmits(['update:modelValue', 'update:bodyType'])

const { t } = useLocale()

const localBodyType = computed({
  get: () => props.bodyType,
  set: (val) => emit('update:bodyType', val),
})

const typeLabel = computed(() => BODY_TYPES.find((o) => o.value === props.bodyType)?.label ?? props.bodyType)

const formRows = ref([{ name: '', value: '', type: 'text', file: null }])

/**
 * Reset editor state
 * @param {string} type - Body type
 * @param {string} text - Text content
 * @param {Array} rows - Form row data
 */
function reset(type, text, rows) {
  emit('update:bodyType', type ?? 'json')
  emit('update:modelValue', text ?? '')
  formRows.value = rows?.length
    ? rows.map((r) => ({ type: 'text', file: null, ...r }))
    : [{ name: '', value: '', type: 'text', file: null }]
}

defineExpose({ formRows, reset })
</script>
