<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { getSvgByProductName, getSvgByType } from '@/utils/svg-assets'

const props = defineProps<{
  src?: string
  alt: string
  className?: string
}>()

const hasError = ref(false)

const isLocal = computed(() => props.src?.startsWith('local:'))

const displaySrc = computed(() => {
  if (isLocal.value && props.src) {
    const parts = props.src.split(':')
    const type = parts[1] ?? 'box'
    return getSvgByType(type)
  }
  return props.src
})

const fallbackUrl = computed(() => getSvgByProductName(props.alt))

// 如果 src 变了，重置错误状态
watch(
  () => props.src,
  () => {
    hasError.value = false
  }
)

const handleError = () => {
  if (!isLocal.value) {
    hasError.value = true
  }
}
</script>

<template>
  <img
    v-if="displaySrc && !hasError"
    :src="displaySrc"
    :alt="alt"
    :class="[className, isLocal ? 'p-8 bg-surface-light object-contain' : '']"
    @error="handleError"
  />
  <img
    v-else
    :src="fallbackUrl"
    :alt="alt"
    :class="[className, 'p-8 bg-surface-light object-contain']"
  />
</template>
