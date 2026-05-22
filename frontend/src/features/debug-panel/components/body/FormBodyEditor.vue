<template>
  <div class="divide-border divide-y">
    <div v-for="(row, i) in formRows" :key="i" class="flex items-center gap-1.5 px-3 py-1.5">
      <!-- Name input -->
      <input
        v-model="row.name"
        :placeholder="t('paramName')"
        class="text-primary-foreground placeholder-muted-foreground min-w-0 flex-1 bg-transparent font-mono text-[11px] focus:outline-none"
      />

      <!-- Type toggle (only for form-data) -->
      <button
        v-if="isFormData"
        @click="toggleRowType(i)"
        class="min-w-8.5 shrink-0 rounded border px-1.5 py-0.5 text-[9px] leading-none font-semibold transition-colors"
        :class="
          row.type === 'file'
            ? 'border-primary/50 text-primary-foreground bg-primary/10'
            : 'border-border text-muted-foreground'
        "
      >
        {{ row.type === 'file' ? 'File' : 'Text' }}
      </button>

      <span class="text-muted-foreground shrink-0 text-[10px]">=</span>

      <!-- File picker -->
      <template v-if="row.type === 'file'">
        <label class="min-w-0 flex-1 cursor-pointer">
          <span
            class="text-muted-foreground block truncate font-mono text-[11px]"
            :class="row.file ? 'text-primary-foreground' : ''"
          >
            {{ row.file ? row.file.name : t('selectFile') }}
          </span>
          <input type="file" class="hidden" @change="(e) => onFileChange(i, e)" />
        </label>
      </template>

      <!-- Text value input -->
      <input
        v-else
        v-model="row.value"
        :placeholder="t('paramValue')"
        class="text-primary-foreground placeholder-muted-foreground min-w-0 flex-1 bg-transparent font-mono text-[11px] focus:outline-none"
      />

      <!-- Delete button -->
      <button @click="removeRow(i)" class="text-muted-foreground shrink-0 transition-colors hover:text-red-500">
        <CircleMinusIcon :size="14" />
      </button>
    </div>

    <!-- Add param button -->
    <div class="px-3 py-2">
      <button
        @click="addRow"
        class="text-muted-foreground hover:bg-accent flex w-full items-center justify-center gap-1 rounded-md py-1.5 text-xs transition-colors"
      >
        <PlusIcon :size="12" />
        {{ t('addParam') }}
      </button>
    </div>
  </div>
</template>

<script setup>
/**
 * FormBodyEditor - Form-based request body editor
 * Supports x-www-form-urlencoded and multipart/form-data
 * FormData mode supports file upload
 * @component
 * @example
 * <FormBodyEditor
 *   v-model="formRows"
 *   :is-form-data="true"
 * />
 */
import { ref, watch } from 'vue'
import { CircleMinusIcon, PlusIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'

const { t } = useLocale()

const props = defineProps({
  /** @type {Array} Form row data */
  modelValue: {
    type: Array,
    required: true,
  },
  /** @type {boolean} Whether in FormData mode (supports file upload) */
  isFormData: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const formRows = ref([...props.modelValue])

function addRow() {
  formRows.value.push({ name: '', value: '', type: 'text', file: null })
  emitUpdate()
}

function removeRow(index) {
  formRows.value.splice(index, 1)
  emitUpdate()
}

function toggleRowType(index) {
  const row = formRows.value[index]
  row.type = row.type === 'file' ? 'text' : 'file'
  if (row.type === 'text') {
    row.file = null
  } else {
    row.value = ''
  }
  emitUpdate()
}

function onFileChange(index, event) {
  const file = event.target.files?.[0]
  if (file) {
    formRows.value[index].file = file
    emitUpdate()
  }
}

function emitUpdate() {
  emit('update:modelValue', formRows.value)
}

watch(
  () => props.modelValue,
  (newVal) => {
    formRows.value = [...newVal]
  },
  { deep: true }
)
</script>
