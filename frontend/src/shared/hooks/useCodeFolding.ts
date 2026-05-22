/**
 * useCodeFolding - Code folding functionality
 * Provides reusable code folding logic for JSON/XML code blocks
 */
import { ref, computed, type Ref, type ComputedRef } from 'vue'
import { detectFoldRanges, getCollapsedIndicator } from '@/shared/utils/codeFolding'

export interface DisplayLine {
  lineNum: number
  content: string
  isFoldStart: boolean
  isCollapsed: boolean
  foldDepth: number
}

export function useCodeFolding(lines: Ref<string[]> | ComputedRef<string[]>, type: 'json' | 'xml' | 'auto' = 'auto') {
  const collapsedLines = ref(new Set<number>())

  function toggleFold(lineNum: number) {
    const next = new Set(collapsedLines.value)
    if (next.has(lineNum)) {
      next.delete(lineNum)
    } else {
      next.add(lineNum)
    }
    collapsedLines.value = next
  }

  function clearFolds() {
    collapsedLines.value.clear()
  }

  const displayLines = computed(() => {
    const ranges = detectFoldRanges(lines.value, type)
    const skipSet = new Set<number>()

    for (const r of ranges) {
      if (collapsedLines.value.has(r.start)) {
        for (let i = r.start + 1; i <= r.end; i++) {
          skipSet.add(i)
        }
      }
    }

    const result: DisplayLine[] = []
    for (let i = 0; i < lines.value.length; i++) {
      if (skipSet.has(i)) continue

      const range = ranges.find((r) => r.start === i)
      const isFoldStart = !!range
      const isCollapsed = isFoldStart && collapsedLines.value.has(i)
      const foldDepth = ranges.filter((r) => r.start < i && r.end > i && !collapsedLines.value.has(r.start)).length

      result.push({
        lineNum: i,
        content: lines.value[i],
        isFoldStart,
        isCollapsed,
        foldDepth,
      })
    }

    return result
  })

  const displayLinesWithCollapsedIndicator = computed(() => {
    return displayLines.value.map((line) => ({
      ...line,
      displayContent: line.isCollapsed
        ? line.content.trimEnd() + getCollapsedIndicator(line.content, type)
        : line.content,
    }))
  })

  return {
    collapsedLines,
    displayLines,
    displayLinesWithCollapsedIndicator,
    toggleFold,
    clearFolds,
  }
}
