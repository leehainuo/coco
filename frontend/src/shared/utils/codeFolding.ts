/**
 * Code folding utilities for JSON and XML
 * Shared logic between CodeBlock and TextBodyEditor components
 */

export interface FoldRange {
  start: number
  end: number
}

/**
 * Detect fold ranges in code lines
 * Supports both JSON and XML folding
 */
export function detectFoldRanges(lines: string[], type: 'json' | 'xml' | 'auto' = 'auto'): FoldRange[] {
  const ranges: FoldRange[] = []
  const jsonStack: number[] = []
  const xmlTagStack: Array<{ tag: string; line: number }> = []

  // Auto-detect type based on content
  let isXml = type === 'xml'
  let isJson = type === 'json'
  if (type === 'auto') {
    isXml = lines.some((line) => line.trim().match(/<\w+[^>]*>/))
    isJson = !isXml
  }

  for (let i = 0; i < lines.length; i++) {
    const trimmed = lines[i].trimEnd()
    const last = trimmed[trimmed.length - 1]

    // JSON folding: { } [ ]
    if (isJson) {
      if (last === '{' || last === '[') {
        jsonStack.push(i)
      }
      const first = lines[i].trim()
      if ((first.startsWith('}') || first.startsWith(']')) && jsonStack.length) {
        const start = jsonStack.pop()!
        if (i - start > 1) {
          ranges.push({ start, end: i })
        }
      }
    }

    // XML folding: <tag> </tag>
    if (isXml) {
      // Check for opening tag (not self-closing and not closed on same line)
      const openTagMatch = trimmed.match(/<(\w+)[^>]*>/)
      if (openTagMatch && !trimmed.includes('/>')) {
        const tagName = openTagMatch[1]
        const closePattern = new RegExp(`</${tagName}>`)
        // Only push to stack if closing tag is not on the same line
        if (!closePattern.test(trimmed)) {
          xmlTagStack.push({ tag: tagName, line: i })
        }
      }

      // Check for closing tag
      const closeTagMatch = trimmed.match(/<\/(\w+)>/)
      if (closeTagMatch && xmlTagStack.length) {
        const tagName = closeTagMatch[1]
        // Find matching opening tag in stack
        for (let j = xmlTagStack.length - 1; j >= 0; j--) {
          if (xmlTagStack[j].tag === tagName) {
            const start = xmlTagStack[j].line
            xmlTagStack.splice(j, 1)
            if (i - start >= 1) {
              ranges.push({ start, end: i })
            }
            break
          }
        }
      }
    }
  }

  return ranges
}

/**
 * Generate collapsed line indicator for a given line
 * Returns the text to append when a line is collapsed
 */
export function getCollapsedIndicator(line: string, type: 'json' | 'xml' | 'auto' = 'auto'): string {
  const trimmed = line.trimEnd()

  // Auto-detect type
  let isXml = type === 'xml'
  if (type === 'auto') {
    isXml = trimmed.includes('<')
  }

  if (isXml) {
    const tagMatch = trimmed.match(/<(\w+)/)
    return ' ⋯ </' + (tagMatch ? tagMatch[1] : 'tag') + '>'
  }

  // JSON
  const closer = trimmed.endsWith('{') ? ' }' : ' ]'
  return ' ⋯' + closer
}
