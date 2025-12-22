<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'

const props = defineProps<{
  targetTime: string
  label?: string
  size?: 'sm' | 'md' | 'lg'
}>()

const emit = defineEmits(['finish'])

const timeLeft = ref(0)
const timer = ref<number | null>(null)

const updateTimer = () => {
  const target = new Date(props.targetTime).getTime()
  const now = Date.now()
  const diff = target - now

  if (diff <= 0) {
    timeLeft.value = 0
    if (timer.value) {
      clearInterval(timer.value)
      timer.value = null
    }
    emit('finish')
  } else {
    timeLeft.value = diff
  }
}

watch(() => props.targetTime, () => {
  updateTimer()
  if (!timer.value && timeLeft.value > 0) {
    timer.value = window.setInterval(updateTimer, 1000)
  }
})

onMounted(() => {
  updateTimer()
  timer.value = window.setInterval(updateTimer, 1000)
})

onUnmounted(() => {
  if (timer.value) clearInterval(timer.value)
})

const duration = computed(() => {
  if (timeLeft.value <= 0) return { hours: '00', minutes: '00', seconds: '00' }
  
  const seconds = Math.floor((timeLeft.value / 1000) % 60)
  const minutes = Math.floor((timeLeft.value / (1000 * 60)) % 60)
  const hours = Math.floor((timeLeft.value / (1000 * 60 * 60)))
  
  return { 
    hours: hours.toString().padStart(2, '0'), 
    minutes: minutes.toString().padStart(2, '0'), 
    seconds: seconds.toString().padStart(2, '0') 
  }
})

const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm': return 'text-sm'
    case 'lg': return 'text-4xl md:text-5xl'
    default: return 'text-xl'
  }
})
</script>

<template>
  <div class="flex flex-col items-start">
    <span v-if="label" class="text-xs text-secondary uppercase tracking-widest mb-1">{{ label }}</span>
    <div :class="['font-mono font-bold tracking-tight flex items-baseline gap-1', sizeClasses]">
      <span class="text-accent">{{ duration.hours }}</span>
      <span class="text-secondary">:</span>
      <span class="text-accent">{{ duration.minutes }}</span>
      <span class="text-secondary">:</span>
      <span class="text-accent">{{ duration.seconds }}</span>
    </div>
  </div>
</template>