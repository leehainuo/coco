<template>
  <div>
    <p class="text-muted-foreground mb-2 text-xs">{{ tag }}</p>

    <div class="mb-4 flex items-center gap-3">
      <h1 class="text-primary-foreground text-2xl leading-tight font-bold">
        {{ title }}
      </h1>
    </div>
  </div>
</template>

<script setup>
defineProps({
  tag: {
    type: String,
    default: '',
  },
  title: {
    type: String,
    required: true,
  },
})
</script>
