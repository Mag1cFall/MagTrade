<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { getSvgByProductName } from '@/utils/svg-assets'

const props = defineProps<{
  src?: string
  alt: string
  className?: string
}>()

const hasError = ref(false)
const svgUrl = computed(() => getSvgByProductName(props.alt))

// 如果 src 变了，重置错误状态
watch(() => props.src, () => {
  hasError.value = false
})

const handleError = () => {
  hasError.value = true
}
</script>

<template>
  <img 
    v-if="src && !hasError" 
    :src="src" 
    :alt="alt" 
    :class="className"
    @error="handleError"
  />
  <img 
    v-else 
    :src="svgUrl" 
    :alt="alt" 
    :class="[className, 'p-8 bg-surface-light object-contain']" 
  />
</template>