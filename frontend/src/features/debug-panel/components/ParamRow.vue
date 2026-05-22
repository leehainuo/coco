<template>
  <div class="flex items-center gap-2">
    <!-- Checkbox -->
    <button
      @click="item.enabled = !item.enabled"
      :class="[
        'text-green-foreground flex h-4 w-4 shrink-0 items-center justify-center rounded border transition-colors',
        item.enabled ? 'border-green-foreground bg-transparent' : 'border-border',
      ]"
    >
      <CheckIcon v-if="item.enabled" :size="10" />
    </button>

    <!-- Name: fixed label or editable input -->
    <span
      v-if="!nameEditable"
      class="text-primary-foreground w-20 shrink-0 cursor-default truncate text-[11px]"
      @mouseenter="item.meta && onEnter($event)"
      @mouseleave="onLeave"
    >
      {{ item.name }}
    </span>
    <input
      v-else
      v-model="item.name"
      :placeholder="t('paramName')"
      class="focus:border-border bg-muted text-muted-foreground w-20 shrink-0 rounded border border-transparent px-2 py-1 font-mono text-[11px] focus:outline-none"
    />

    <!-- Value -->
    <input
      v-model="item.value"
      :placeholder="t('paramValue')"
      class="focus:border-border bg-muted text-muted-foreground min-w-0 flex-1 rounded border border-transparent px-2 py-1 font-mono text-[11px] focus:outline-none"
    />

    <!-- Delete -->
    <button @click="$emit('remove')" class="text-muted-foreground hover:text-destructive shrink-0 transition-colors">
      <CircleMinusIcon :size="14" />
    </button>

    <!-- Tooltip teleported to body so it escapes overflow-hidden -->
    <Teleport to="body">
      <div
        v-if="tooltip.visible"
        :style="{ top: tooltip.top + 'px', left: tooltip.left + 'px' }"
        class="bg-background border-border pointer-events-none fixed z-9999 max-w-[220px] min-w-[160px] overflow-hidden rounded-lg border text-[11px] shadow-xl"
      >
        <div class="border-border grid grid-cols-[auto_1fr] items-center gap-x-4 border-b px-3 py-1.5">
          <span class="text-muted-foreground whitespace-nowrap">{{ t('paramType') }}</span>
          <span class="text-foreground truncate font-mono">{{ item.meta?.type ?? '-' }}</span>
        </div>
        <div class="border-border grid grid-cols-[auto_1fr] items-center gap-x-4 border-b px-3 py-1.5">
          <span class="text-muted-foreground whitespace-nowrap">{{ t('paramRequired') }}</span>
          <span :class="item.meta?.required ? 'text-destructive' : 'text-muted-foreground'">
            {{ item.meta?.required ? '✓' : '✗' }}
          </span>
        </div>
        <div class="grid grid-cols-[auto_1fr] items-start gap-x-4 px-3 py-1.5">
          <span class="text-muted-foreground whitespace-nowrap">{{ t('paramDesc') }}</span>
          <span class="text-foreground leading-relaxed">{{ item.meta?.description || '-' }}</span>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useLocale } from '@/shared/hooks'
import { CheckIcon, CircleMinusIcon } from '@/shared/components/icons'

defineProps({
  item: { type: Object, required: true },
  nameEditable: { type: Boolean, default: false },
})

defineEmits(['remove'])

const { t } = useLocale()

const tooltip = ref({ visible: false, top: 0, left: 0 })

function onEnter(event) {
  const rect = event.currentTarget.getBoundingClientRect()
  tooltip.value = { visible: true, top: rect.top, left: rect.right - 60 }
}

function onLeave() {
  tooltip.value.visible = false
}
</script>
