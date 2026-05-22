<template>
  <div class="border-border mx-3 my-2 overflow-hidden rounded-lg border">
    <!-- Header -->
    <div class="border-border flex shrink-0 items-center justify-between border-b px-3 py-1">
      <span class="text-muted-foreground text-[10px]">{{ typeLabel }}</span>
      <button @click="copy" class="text-muted-foreground hover:text-accent-foreground transition-colors">
        <CheckIcon v-if="copied" :size="12" />
        <CopyIcon v-else :size="12" />
      </button>
    </div>

    <!-- Editor area -->
    <div class="overflow-auto" style="max-height: 130px; min-height: 130px">
      <div class="flex">
        <!-- Gutter: fold controls + line numbers + scope lines -->
        <div class="flex shrink-0 flex-col pt-3 pb-3 select-none">
          <div
            v-for="line in displayLines"
            :key="line.lineNum"
            class="text-muted-foreground flex h-5 items-center font-mono text-[10px]"
          >
            <button
              v-if="line.isFoldStart"
              @click.stop="toggleFold(line.lineNum)"
              class="flex w-4 shrink-0 items-center justify-center opacity-60 transition-opacity hover:opacity-100"
            >
              <ChevronDownIcon v-if="!collapsedLines.has(line.lineNum)" :size="9" />
              <ChevronRightIcon v-else :size="9" />
            </button>
            <span v-else class="w-4 shrink-0" />
            <span class="shrink-0 pr-1 text-right" style="min-width: 1.25rem">{{ line.lineNum + 1 }}</span>
            <div class="relative flex w-3 shrink-0 justify-center self-stretch">
              <div v-if="line.foldDepth > 0" class="bg-border h-full w-px" />
            </div>
          </div>
        </div>

        <!-- Editor: highlight pre + transparent textarea -->
        <div class="relative min-w-0 flex-1">
          <pre
            class="text-primary-foreground pointer-events-none w-full px-2 py-3 font-mono text-[11px] leading-5 break-all whitespace-pre-wrap"
            v-html="visibleHighlighted"
          />
          <textarea
            :value="modelValue"
            @input="$emit('update:modelValue', $event.target.value)"
            @keydown.tab.prevent="onTab"
            spellcheck="false"
            class="caret-primary-foreground absolute inset-0 h-full w-full resize-none overflow-hidden bg-transparent px-2 py-3 font-mono text-[11px] leading-5 text-transparent focus:outline-none"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * TextBodyEditor - Text-based request body editor
 * Supports JSON, XML, Text, Binary, GraphQL and other text types
 * Features: syntax highlighting, code folding, line numbers, copy
 * @component
 * @example
 * <TextBodyEditor
 *   v-model="bodyText"
 *   body-type="json"
 *   type-label="JSON"
 * />
 */
import { ref, computed, watch } from 'vue'
import { CheckIcon, CopyIcon, ChevronDownIcon, ChevronRightIcon } from '@/shared/components/icons'
import { highlightJson, highlightXml } from '@/shared/utils/highlight'
import { detectFoldRanges } from '@/shared/utils/codeFolding'
import { useClipboard } from '@/shared/hooks'

const props = defineProps({
  /** @type {string} Text content */
  modelValue: {
    type: String,
    required: true,
  },
  /** @type {string} Request body type */
  bodyType: {
    type: String,
    required: true,
  },
  /** @type {string} Type display label */
  typeLabel: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['update:modelValue'])

const { copied, copy: copyToClipboard } = useClipboard(2000)
const collapsedLines = ref(new Set())

const lines = computed(() => props.modelValue.split('\n'))

const displayLines = computed(() => {
  const foldType = props.bodyType === 'xml' ? 'xml' : 'json'
  const ranges = detectFoldRanges(lines.value, foldType)

  const result = []
  const stack = []

  for (let i = 0; i < lines.value.length; i++) {
    const text = lines.value[i]
    const range = ranges.find((r) => r.start === i)
    const isFoldStart = !!range

    // Check if this line should be hidden (inside a collapsed range)
    const isCollapsed = stack.some((lineNum) => collapsedLines.value.has(lineNum))

    // Check if we're at the end of a fold range
    if (ranges.some((r) => r.end === i) && stack.length > 0) {
      stack.pop()
    }

    if (!isCollapsed) {
      result.push({
        lineNum: i,
        text,
        isFoldStart,
        foldDepth: stack.length,
      })
    }

    if (isFoldStart) {
      stack.push(i)
    }
  }

  return result
})

const visibleHighlighted = computed(() => {
  const visibleText = displayLines.value.map((l) => l.text).join('\n')
  if (props.bodyType === 'json') {
    return highlightJson(visibleText)
  } else if (props.bodyType === 'xml') {
    return highlightXml(visibleText)
  }
  // For other types (text, binary, graphql), escape HTML to display correctly
  return visibleText.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
})

function toggleFold(lineNum) {
  const next = new Set(collapsedLines.value)
  next.has(lineNum) ? next.delete(lineNum) : next.add(lineNum)
  collapsedLines.value = next
}

function onTab(e) {
  const start = e.target.selectionStart
  const end = e.target.selectionEnd
  const value = e.target.value
  const newValue = value.substring(0, start) + '  ' + value.substring(end)
  emit('update:modelValue', newValue)
  setTimeout(() => {
    e.target.selectionStart = e.target.selectionEnd = start + 2
  }, 0)
}

function copy() {
  copyToClipboard(props.modelValue)
}

watch(
  () => props.modelValue,
  () => {
    collapsedLines.value.clear()
  }
)
</script>
