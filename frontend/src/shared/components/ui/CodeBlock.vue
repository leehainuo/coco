<template>
  <div class="font-mono text-xs leading-relaxed">
    <div v-for="line in displayLinesWithHL" :key="line.lineNum" class="flex">
      <!-- Gutter: chevron + line number + scope line -->
      <div class="text-muted-foreground flex shrink-0 items-start text-[10px] select-none">
        <button
          v-if="line.isFoldStart"
          @click.stop="toggleFold(line.lineNum)"
          class="mt-[0.2em] flex w-4 shrink-0 justify-center opacity-60 transition-opacity hover:opacity-100"
        >
          <ChevronDownIcon v-if="!collapsedLines.has(line.lineNum)" :size="9" />
          <ChevronRightIcon v-else :size="9" />
        </button>
        <span v-else class="w-4 shrink-0" />
        <span class="mt-[0.1em] shrink-0 pr-1 text-right" style="min-width: 1.25rem">{{ line.lineNum + 1 }}</span>
        <div class="relative flex w-3 shrink-0 justify-center self-stretch">
          <div v-if="line.foldDepth > 0" class="bg-border h-full w-px" />
        </div>
      </div>
      <!-- Content -->
      <pre class="text-primary-foreground min-w-0 flex-1 break-all whitespace-pre-wrap" v-html="line.highlighted" />
    </div>
  </div>
</template>

<script setup>
/**
 * CodeBlock - Code block component
 * Supports JSON and XML syntax highlighting, code folding, line numbers
 * @component
 * @example
 * <CodeBlock :code="jsonString" />
 * <CodeBlock :code="xmlString" />
 */
import { computed } from 'vue'
import { ChevronDownIcon, ChevronRightIcon } from '@/shared/components/icons'
import { highlightCode } from '@/shared/utils/highlight'
import { useCodeFolding } from '@/shared/hooks'

const props = defineProps({
  /** @type {string} Code content to display (JSON or XML format) */
  code: { type: String, default: '' },
})

const lines = computed(() => (props.code || '').split('\n'))
const { collapsedLines, displayLinesWithCollapsedIndicator, toggleFold } = useCodeFolding(lines, 'auto')

const displayLinesWithHL = computed(() => {
  const code = displayLinesWithCollapsedIndicator.value.map((line) => line.displayContent).join('\n')
  const hlLines = highlightCode(code).split('\n')
  return displayLinesWithCollapsedIndicator.value.map((line, idx) => ({
    ...line,
    highlighted: hlLines[idx] ?? '',
  }))
})
</script>
