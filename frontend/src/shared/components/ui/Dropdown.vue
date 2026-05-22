<template>
  <div class="relative inline-block">
    <!-- Trigger slot -->
    <div ref="triggerRef" @click="toggle">
      <slot name="trigger" :open="isOpen" />
    </div>

    <!-- Dropdown menu via Teleport -->
    <Teleport to="body">
      <!-- Backdrop -->
      <div v-if="isOpen" class="fixed inset-0 z-9998" @click="close" />

      <!-- Dropdown content -->
      <div
        v-if="isOpen"
        ref="dropdownRef"
        :style="dropdownStyle"
        :class="['bg-background border-border fixed z-9999 overflow-hidden rounded-lg border shadow-xl', contentClass]"
      >
        <slot :close="close" />
      </div>
    </Teleport>
  </div>
</template>

<script setup>
/**
 * Dropdown - Generic dropdown menu component
 * Supports 4 positions, auto viewport overflow handling, ESC key to close
 * @component
 * @example
 * <Dropdown placement="bottom-end">
 *   <template #trigger="{ open }">
 *     <button>Open Menu</button>
 *   </template>
 *   <template #default="{ close }">
 *     <div>Menu Content</div>
 *   </template>
 * </Dropdown>
 */
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'

const props = defineProps({
  /** @type {string} Dropdown menu position */
  placement: {
    type: String,
    default: 'bottom-start',
    validator: (v) => ['bottom-start', 'bottom-end', 'top-start', 'top-end'].includes(v),
  },
  /** @type {number} Spacing from trigger (pixels) */
  offset: {
    type: Number,
    default: 4,
  },
  /** @type {string} Additional CSS class for content area */
  contentClass: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['open', 'close'])

const isOpen = ref(false)
const triggerRef = ref(null)
const dropdownRef = ref(null)
const dropdownStyle = ref({})

function toggle() {
  if (isOpen.value) {
    close()
  } else {
    open()
  }
}

async function open() {
  isOpen.value = true
  emit('open')
  await nextTick()
  updatePosition()
}

function close() {
  isOpen.value = false
  emit('close')
}

function updatePosition() {
  if (!triggerRef.value || !dropdownRef.value) return

  const triggerRect = triggerRef.value.getBoundingClientRect()
  const dropdownRect = dropdownRef.value.getBoundingClientRect()
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight

  let top = 0
  let left = 0

  // Calculate position based on placement
  switch (props.placement) {
    case 'bottom-start':
      top = triggerRect.bottom + props.offset
      left = triggerRect.left
      break
    case 'bottom-end':
      top = triggerRect.bottom + props.offset
      left = triggerRect.right - dropdownRect.width
      break
    case 'top-start':
      top = triggerRect.top - dropdownRect.height - props.offset
      left = triggerRect.left
      break
    case 'top-end':
      top = triggerRect.top - dropdownRect.height - props.offset
      left = triggerRect.right - dropdownRect.width
      break
  }

  // Adjust if overflow viewport
  if (left + dropdownRect.width > viewportWidth) {
    left = viewportWidth - dropdownRect.width - 8
  }
  if (left < 8) {
    left = 8
  }
  if (top + dropdownRect.height > viewportHeight) {
    top = triggerRect.top - dropdownRect.height - props.offset
  }
  if (top < 8) {
    top = triggerRect.bottom + props.offset
  }

  dropdownStyle.value = {
    top: `${top}px`,
    left: `${left}px`,
  }
}

// Close on escape key
function handleKeydown(e) {
  if (e.key === 'Escape' && isOpen.value) {
    close()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('resize', updatePosition)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('resize', updatePosition)
})

defineExpose({ open, close, toggle })
</script>
