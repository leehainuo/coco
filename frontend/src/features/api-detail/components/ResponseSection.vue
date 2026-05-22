<template>
  <div class="space-y-2">
    <div v-for="(resp, code) in responses" :key="code" class="border-border overflow-hidden rounded-lg border">
      <button
        class="hover:bg-accent flex w-full items-center gap-3 px-4 py-3 text-left transition-colors"
        @click="toggleResp(code)"
      >
        <StatusDot :code="code" :description="resp.description" />
        <ChevronDownIcon
          v-if="expandedResps.has(code)"
          :size="14"
          color="var(--icon-foreground)"
          class="ml-auto shrink-0"
        />
        <ChevronRightIcon v-else :size="14" color="var(--icon-foreground)" class="ml-auto shrink-0" />
      </button>
      <div v-show="expandedResps.has(code)" class="border-border border-t px-4 py-4">
        <!-- OpenAPI 3.0 format (content) -->
        <template v-if="resp.content">
          <div v-for="(mediaObj, mediaType) in resp.content" :key="mediaType">
            <div class="grid grid-cols-2 gap-6">
              <div>
                <span class="text-base-foreground font-mono text-xs">{{ mediaType }}</span>
                <p class="mt-2 text-base text-sm wrap-break-word">{{ resolveType(mediaObj.schema) }}</p>
              </div>
              <div>
                <p class="text-muted-foreground mb-2 text-xs">{{ t('example') }}</p>
                <CodeBlock :code="formatSchema(mediaObj.schema, mediaType)" />
              </div>
            </div>
          </div>
        </template>
        <!-- Swagger 2.0 format (schema) -->
        <template v-else-if="resp.schema">
          <div class="grid grid-cols-2 gap-6">
            <div>
              <span class="text-base-foreground font-mono text-xs">{{ produces }}</span>
              <p class="mt-2 text-base text-sm wrap-break-word">{{ resolveType(resp.schema) }}</p>
            </div>
            <div>
              <p class="text-muted-foreground mb-2 text-xs">{{ t('example') }}</p>
              <CodeBlock :code="formatSchema(resp.schema, produces)" />
            </div>
          </div>
        </template>
        <p v-else class="text-muted-foreground text-xs">{{ t('noResponseBody') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { StatusDot, CodeBlock } from '@/shared/components/ui'
import { ChevronDownIcon, ChevronRightIcon } from '@/shared/components/icons'
import { resolveType, makeSchemaHelpers } from '@/shared/utils/schema'
import { useLocale } from '@/shared/hooks'

const props = defineProps({
  responses: { type: Object, required: true },
  spec: { type: Object, default: null },
  produces: { type: String, default: 'application/json' },
})

function formatSchema(schema, mediaType) {
  return makeSchemaHelpers(props.spec).formatSchema(schema, mediaType)
}

const { t } = useLocale()
const expandedResps = ref(new Set())

watch(
  () => props.responses,
  () => {
    expandedResps.value = new Set()
  }
)

function toggleResp(code) {
  const s = new Set(expandedResps.value)
  s.has(code) ? s.delete(code) : s.add(code)
  expandedResps.value = s
}
</script>
