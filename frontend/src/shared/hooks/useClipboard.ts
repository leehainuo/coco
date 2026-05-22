/**
 * useClipboard - Clipboard copy functionality
 * Provides a reusable way to copy text to clipboard with visual feedback
 */
import { ref } from 'vue'

export function useClipboard(duration = 1500) {
  const copied = ref(false)

  async function copy(text: string) {
    try {
      await navigator.clipboard.writeText(text)
      copied.value = true
      setTimeout(() => {
        copied.value = false
      }, duration)
      return true
    } catch (error) {
      console.error('Failed to copy to clipboard:', error)
      return false
    }
  }

  return {
    copied,
    copy,
  }
}
