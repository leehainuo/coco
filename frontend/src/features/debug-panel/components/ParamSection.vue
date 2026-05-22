<template>
  <section class="shrink-0">
    <p class="text-primary-foreground px-3 pt-2.5 pb-1.5 text-xs font-medium">{{ title }}</p>
    <div class="space-y-1.5 px-3 pb-3">
      <ParamRow
        v-for="(item, i) in items"
        :key="i"
        :item="item"
        :name-editable="item.custom ?? nameEditable"
        @remove="$emit('remove', i)"
      />
      <button
        v-if="addable"
        @click="$emit('add')"
        class="text-muted-foreground hover:text-accent-foreground hover:bg-accent inline-flex items-center gap-1 rounded px-2 py-1 text-xs text-[11px] transition-colors"
      >
        <PlusIcon :size="12" class="shrink-0" />
        {{ t('addParam') }}
      </button>
    </div>
  </section>
</template>

<script setup>
import { ParamRow } from '@/features/debug-panel'
import { PlusIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'

defineProps({
  title: { type: String, required: true },
  items: { type: Array, required: true },
  addable: { type: Boolean, default: false },
  nameEditable: { type: Boolean, default: false },
})

defineEmits(['add', 'remove'])

const { t } = useLocale()
</script>
