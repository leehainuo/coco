<template>
  <div class="border-border overflow-x-auto rounded-md border">
    <table class="w-full min-w-max text-xs">
      <thead class="bg-muted">
        <tr>
          <th class="text-muted-foreground px-4 py-2 text-left font-medium whitespace-nowrap">
            {{ t('paramName') }}
          </th>
          <th class="text-muted-foreground px-4 py-2 text-left font-medium whitespace-nowrap">
            {{ t('paramType') }}
          </th>
          <th class="text-muted-foreground px-4 py-2 text-left font-medium whitespace-nowrap">
            {{ t('paramRequired') }}
          </th>
          <th class="text-muted-foreground px-4 py-2 text-left font-medium whitespace-nowrap">
            {{ t('paramDesc') }}
          </th>
        </tr>
      </thead>
      <tbody class="divide-border divide-y">
        <tr v-for="p in params" :key="p.name">
          <td class="text-primary-foreground px-4 py-2.5 font-mono whitespace-nowrap">{{ p.name }}</td>
          <td class="px-4 py-2.5 font-mono text-base whitespace-nowrap">
            {{ getParamType(p) }}
          </td>
          <td class="px-4 py-2.5 whitespace-nowrap">
            <span v-if="p.required" class="text-destructive font-semibold">{{ t('yes') }}</span>
            <span v-else class="text-muted-foreground">{{ t('no') }}</span>
          </td>
          <td class="text-primary-foreground px-4 py-2.5 whitespace-nowrap">{{ p.description || '—' }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { resolveType } from '@/shared/utils/schema'
import { useLocale } from '@/shared/hooks'

defineProps({
  params: { type: Array, default: () => [] },
})

const { t } = useLocale()

/**
 * Get parameter type - handles both Swagger 2.0 and OpenAPI 3.0 formats
 * Swagger 2.0: type is directly on the parameter object
 * OpenAPI 3.0: type is in the schema object
 */
function getParamType(param) {
  // Swagger 2.0 style: type directly on parameter
  if (param.type) {
    if (param.type === 'array' && param.items) {
      const itemType = param.items.type || param.items.$ref?.split('/').pop() || 'any'
      return `array[${itemType}]`
    }
    return param.type
  }
  // OpenAPI 3.0 style: type in schema
  return resolveType(param.schema)
}
</script>
