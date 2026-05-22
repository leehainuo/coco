<template>
  <section v-if="responses" class="mb-8">
    <h2 class="text-primary-foreground mb-4 font-semibold">{{ t('response') }}</h2>
    <ResponseSection :responses="responses" :spec="spec" :produces="produces" />
  </section>
</template>

<script setup>
import { ResponseSection } from '@/features/api-detail'
import { useLocale } from '@/shared/hooks'

const { t } = useLocale()

defineProps({
  responses: {
    type: Object,
    default: null,
  },
  spec: {
    type: Object,
    default: null,
  },
  produces: {
    type: String,
    default: 'application/json',
  },
})
</script>
