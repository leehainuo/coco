<template>
  <!-- OpenAPI 3.0 format (requestBody.content) -->
  <div
    v-if="requestBody.content"
    v-for="(mediaObj, mediaType) in requestBody.content"
    :key="mediaType"
    class="@container mb-4"
  >
    <div class="mb-3 flex items-center gap-2">
      <span class="text-primary-foreground text-sm font-medium">{{ t('bodyParam') }}</span>
      <Badge size="xs">{{ mediaType }}</Badge>
      <Badge v-if="requestBody.required" variant="red" size="xs">{{ t('required') }}</Badge>
    </div>
    <div class="flex flex-col gap-6 @lg:grid @lg:grid-cols-2">
      <div class="pt-1 text-base text-sm wrap-break-word">
        {{ resolveType(mediaObj.schema) }}
      </div>
      <div>
        <p class="text-muted-foreground mb-2 text-xs">{{ t('example') }}</p>
        <CodeBlock :code="formatSchema(mediaObj.schema, mediaType)" />
      </div>
    </div>
  </div>

  <!-- Swagger 2.0 format (requestBody.schema) -->
  <div v-else-if="requestBody.schema" class="@container mb-4">
    <div class="mb-3 flex items-center gap-2">
      <span class="text-primary-foreground text-sm font-medium">{{ t('bodyParam') }}</span>
      <Badge size="xs">{{ contentType }}</Badge>
      <Badge v-if="requestBody.required" variant="red" size="xs">{{ t('required') }}</Badge>
    </div>
    <div class="flex flex-col gap-6 @lg:grid @lg:grid-cols-2">
      <div class="pt-1 text-base text-sm wrap-break-word">
        {{ resolveType(requestBody.schema) }}
      </div>
      <div>
        <p class="text-muted-foreground mb-2 text-xs">{{ t('example') }}</p>
        <CodeBlock :code="formatSchema(requestBody.schema, contentType)" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Badge, CodeBlock } from '@/shared/components/ui'
import { resolveType, makeSchemaHelpers } from '@/shared/utils/schema'
import { useLocale } from '@/shared/hooks'

const props = defineProps({
  requestBody: { type: Object, required: true },
  spec: { type: Object, default: null },
  consumes: { type: Array, default: () => [] },
})

// Get Content-Type for Swagger 2.0 (from consumes array)
const contentType = computed(() => {
  if (props.consumes && props.consumes.length > 0) {
    return props.consumes[0]
  }
  return 'application/json'
})

function formatSchema(schema, mediaType) {
  return makeSchemaHelpers(props.spec).formatSchema(schema, mediaType)
}

const { t } = useLocale()
</script>
