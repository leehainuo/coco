<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200"
      leave-active-class="transition-opacity duration-300"
      enter-from-class="opacity-0"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isLoading"
        class="fixed top-0 left-0 z-9999 h-1 bg-linear-to-r from-blue-500 via-purple-500 to-pink-500"
        :style="{ width: progress + '%' }"
      />
    </Transition>
  </Teleport>
</template>

<script setup>
/**
 * TopLoadingBar - Top loading progress bar
 * A horizontal loading bar that appears at the top of the page
 * @component
 * @example
 * <TopLoadingBar :loading="isLoading" />
 */
import { ref, watch, onUnmounted } from 'vue'

const props = defineProps({
  /** @type {boolean} Whether loading is active */
  loading: { type: Boolean, default: false },
  /** @type {number} Animation duration in milliseconds */
  duration: { type: Number, default: 2000 },
})

const isLoading = ref(false)
const progress = ref(0)
let timer = null
let startTime = 0

watch(
  () => props.loading,
  (newVal) => {
    if (newVal) {
      start()
    } else {
      finish()
    }
  },
  { immediate: true }
)

function start() {
  isLoading.value = true
  progress.value = 0
  startTime = Date.now()
  animate()
}

function animate() {
  if (!isLoading.value) return

  const elapsed = Date.now() - startTime
  const percentage = Math.min((elapsed / props.duration) * 100, 95)
  progress.value = percentage

  if (percentage < 95) {
    timer = requestAnimationFrame(animate)
  }
}

function finish() {
  progress.value = 100
  setTimeout(() => {
    isLoading.value = false
    progress.value = 0
    if (timer) {
      cancelAnimationFrame(timer)
      timer = null
    }
  }, 300)
}

onUnmounted(() => {
  if (timer) {
    cancelAnimationFrame(timer)
  }
})
</script>
