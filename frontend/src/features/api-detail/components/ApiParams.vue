<template>
  <section v-if="hasParams" class="mb-8">
    <h2 class="text-primary-foreground mb-4 font-semibold">{{ t('requestParams') }}</h2>

    <!-- Parameters by location -->
    <template v-if="params.length">
      <div v-for="loc in paramLocations" :key="loc" class="mb-5">
        <div class="mb-3 flex items-center gap-2">
          <span class="text-primary-foreground text-sm font-medium capitalize">{{ paramLocLabel(loc) }}</span>
          <Badge v-for="badge in paramBadges(loc)" :key="badge" size="xs">{{ badge }}</Badge>
        </div>
        <ParamTable :params="paramsByLoc(loc)" />
      </div>
    </template>

    <!-- Request body (OpenAPI 3.0 and Swagger 2.0) -->
    <RequestBodySection v-if="requestBody" :request-body="requestBody" :spec="spec" :consumes="consumes" />
  </section>
</template>

<script setup>
import { computed } from 'vue'
import { Badge } from '@/shared/components/ui'
import { ParamTable, RequestBodySection } from '@/features/api-detail'
import { useLocale } from '@/shared/hooks'
import { PARAM_LOCATIONS } from '@/shared/constants/paramLocations'

const { t } = useLocale()

const props = defineProps({
  params: {
    type: Array,
    default: () => [],
  },
  requestBody: {
    type: Object,
    default: null,
  },
  spec: {
    type: Object,
    default: null,
  },
  consumes: {
    type: Array,
    default: () => [],
  },
})

const hasParams = computed(() => props.params.length > 0 || props.requestBody)

const paramLocations = computed(() => {
  const locs = new Set()
  props.params.forEach((p) => {
    if (p.in) locs.add(p.in)
  })
  return PARAM_LOCATIONS.filter((loc) => locs.has(loc))
})

function paramsByLoc(loc) {
  return props.params.filter((p) => p.in === loc)
}

function paramLocLabel(loc) {
  const labels = {
    path: 'Path',
    query: 'Query',
    header: 'Header',
    cookie: 'Cookie',
    formData: 'Form Data',
    body: 'Body',
  }
  return labels[loc] || loc
}

function paramBadges(loc) {
  const params = paramsByLoc(loc)
  const badges = []
  const requiredCount = params.filter((p) => p.required).length
  if (requiredCount > 0) {
    badges.push(`${requiredCount} ${t('required')}`)
  }
  return badges
}
</script>
