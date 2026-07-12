/**
 * Internationalization composable
 */

import { ref, computed } from 'vue'
import type { Locale } from '../types'
import { messages, customLanguageName, registerCustomMessages, type MessageKey } from '../utils/i18n'

const STORAGE_KEY = 'coco:lang'

const saved = localStorage.getItem(STORAGE_KEY) as Locale | null
const locale = ref<Locale>(saved && messages[saved] ? saved : 'zh')

/** Whether a custom language has been registered from i18n.json. */
const hasCustom = ref(false)

export function useLocale() {
  /**
   * Translation function. Falls back to English, then to the raw key.
   */
  function t(key: MessageKey | string): string {
    return messages[locale.value]?.[key] ?? messages.en[key] ?? key
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

  /**
   * Register a custom language from a parsed i18n.json payload and expose
   * it as a selectable locale.
   */
  function loadCustomMessages(payload: unknown): boolean {
    const ok = registerCustomMessages(payload)
    hasCustom.value = ok
    return ok
  }

  /** Display name of the custom language, for the language switcher. */
  const customName = computed(() => customLanguageName.value)

  return {
    locale,
    t,
    setLocale,
    loadCustomMessages,
    hasCustom,
    customName,
  }
}
