/**
 * Internationalization composable
 */

import { ref } from 'vue'
import type { Locale } from '../types'
import { messages, type MessageKey } from '../utils/i18n'

const STORAGE_KEY = 'coco:lang'

const saved = localStorage.getItem(STORAGE_KEY) as Locale | null
const locale = ref<Locale>(saved && messages[saved] ? saved : 'zh')

export function useLocale() {
  /**
   * Translation function
   */
  function t(key: MessageKey | string): string {
    return messages[locale.value]?.[key] ?? messages.zh[key] ?? key
  }

  /**
   * Set locale
   */
  function setLocale(lang: Locale, persist = true): void {
    if (messages[lang]) {
      locale.value = lang
      if (persist) {
        localStorage.setItem(STORAGE_KEY, lang)
      }
    }
  }

  return {
    locale,
    t,
    setLocale,
  }
}
