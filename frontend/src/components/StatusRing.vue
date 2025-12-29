<template>
  <div class="status-ring">
    <svg viewBox="0 0 100 100" class="ring-svg">
      <!-- Track -->
      <circle cx="50" cy="50" r="45" fill="none" stroke="rgba(255,255,255,0.1)" stroke-width="6" />

      <!-- Progress Indicator -->
      <circle
        cx="50"
        cy="50"
        r="45"
        fill="none"
        stroke="#e33535"
        stroke-width="6"
        stroke-linecap="round"
        stroke-dasharray="283"
        :stroke-dashoffset="dashOffset"
        class="progress-circle"
      />
    </svg>

    <div class="ring-content">
      <div class="value">{{ value }}%</div>
      <div class="label">{{ label }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  value: number
  label: string
}>()

const dashOffset = computed(() => {
  const circumference = 2 * Math.PI * 45 // 283
  return circumference - (props.value / 100) * circumference
})
</script>

<style scoped>
.status-ring {
  position: relative;
  width: 120px;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ring-svg {
  transform: rotate(-90deg);
  width: 100%;
  height: 100%;
}

.progress-circle {
  transition: stroke-dashoffset 1s cubic-bezier(0.4, 0, 0.2, 1);
  filter: drop-shadow(0 0 4px rgba(227, 53, 53, 0.5));
}

.ring-content {
  position: absolute;
  text-align: center;
}

.value {
  font-size: 24px;
  font-weight: bold;
  font-family: 'JetBrains Mono', monospace;
  color: #fff;
}

.label {
  font-size: 10px;
  color: #666;
  text-transform: uppercase;
  margin-top: 2px;
}
</style>
