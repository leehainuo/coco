<template>
  <div class="border-border mt-auto flex items-center justify-between border-t pt-4">
    <div v-if="enableExport" class="flex items-center gap-4">
      <Button as="a" variant="ghost" size="sm" :href="openapiJsonPath" download="openapi.json">
        <DownloadIcon :size="14" class="shrink-0" />
        {{ t('export') }}
      </Button>
    </div>
    <span class="text-muted-foreground text-xs font-bold">Built with Coco.</span>
  </div>
</template>

<script setup>
import { Button } from '@/shared/components/ui'
import { DownloadIcon } from '@/shared/components/icons'
import { useLocale } from '@/shared/hooks'

defineProps({
  openapiJsonPath: { type: String, default: '' },
  enableExport: { type: Boolean, default: true },
})

const { t } = useLocale()
</script>
