/**
 * Highlight JSON syntax with Tailwind CSS classes
 * Matches VS Code style: blue keys, emerald values, violet booleans
 */
export function highlightJson(json: string, _multiline = false): string {
  const escapeHtml = (str: string) =>
    str.replace(/[&<>"']/g, (m) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[m] || m)

  // First escape the entire string
  const escaped = escapeHtml(json)

  // Then apply syntax highlighting to JSON tokens
  return escaped.replace(
    /(&quot;(?:\\u[\dA-Fa-f]{4}|\\[^u]|[^&])*?&quot;(?:\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)/g,
    (m) => {
      if (m.startsWith('&quot;')) {
        const isKey = /:$/.test(m.trimEnd())
        return `<span class="${isKey ? 'text-blue-600 dark:text-blue-400' : 'text-emerald-600 dark:text-emerald-400'}">${m}</span>`
      }
      if (m === 'true' || m === 'false') return `<span class="text-violet-600 dark:text-violet-400">${m}</span>`
      if (m === 'null') return `<span class="text-gray-400">${m}</span>`
      return `<span class="text-amber-600 dark:text-amber-400">${m}</span>`
    }
  )
}

/**
 * Highlight XML with custom colors
 * Brackets: #ba3c98, Tags: #f7768e, Content: #c2cad1
 */
export function highlightXml(xml: string): string {
  // First escape HTML
  let result = xml.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')

  // Then apply syntax highlighting with custom colors
  result = result
    // Opening/closing tags: brackets (#ba3c98) + tag names (#f7768e)
    .replace(/(&lt;\/?)(\w+)/g, '<span style="color: #ba3c98">$1</span><span style="color: #f7768e">$2</span>')
    // Closing bracket > (#ba3c98)
    .replace(/(&gt;)/g, '<span style="color: #ba3c98">$1</span>')
    // Attributes: name="value"
    .replace(
      /(\w+)=(&quot;[^&]*?&quot;)/g,
      '<span style="color: #f7768e">$1=</span><span style="color: #c2cad1">$2</span>'
    )
    // Text content between tags (#c2cad1)
    .replace(/(&gt;)([^&<]+?)(&lt;)/g, '$1<span style="color: #c2cad1">$2</span>$3')

  return result
}

/**
 * Auto-detect and highlight code (JSON or XML)
 */
export function highlightCode(code: string): string {
  if (!code) return ''
  const trimmed = code.trim()
  if (trimmed.startsWith('<') && trimmed.includes('</')) {
    return highlightXml(code)
  }
  return highlightJson(code)
}
